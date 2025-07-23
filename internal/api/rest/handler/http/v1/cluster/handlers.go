package cluster

import (
	"gallery-service/config"
	"gallery-service/internal/api/rest/validator"
	clusterCommands "gallery-service/internal/application/commands/v1/cluster"
	requests "gallery-service/internal/application/dto/requests/cluster"
	"gallery-service/internal/application/dto/responses"
	clusterQueries "gallery-service/internal/application/queries/cluster"
	"gallery-service/internal/domain/service"
	constants2 "gallery-service/internal/pkg/constants"
	"gallery-service/pkg/constants"
	httpPkg "gallery-service/pkg/http"
	"gallery-service/pkg/utils"
	"gallery-service/pkg/zap"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type clusterHandlers struct {
	log         zap.Logger
	cfg         *config.Config
	ps          *service.ClusterService
	val         *validator.Wrapper
	mongoClient *mongo.Client
}

func NewClusterHandlers(
	log zap.Logger,
	cfg *config.Config,
	mongoClient *mongo.Client,
) *clusterHandlers {
	return &clusterHandlers{
		log:         log,
		cfg:         cfg,
		val:         validator.NewValidator(log, cfg),
		mongoClient: mongoClient,
	}
}

// CreateCluster
// @Tags clusters
// @Summary Create Cluster
// @Description Create new Cluster
// @Param Cluster body dto.CreateClusterReqDto true "create Cluster"
// @Accept json
// @Produce json
// @Success 201 {string} id ""
// @Router /clusters [post]
func (p *clusterHandlers) CreateCluster(c *fiber.Ctx) error {
	ctx := c.Context()
	var reqDto requests.CreateClusterReqDto
	if err := c.BodyParser(&reqDto); err != nil {
		p.log.Errorf("(Bind) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}
	err := p.val.DataValidation(reqDto)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	command := clusterCommands.NewCreateClusterCommand(
		reqDto.ClusterName,
		reqDto.Title,
		reqDto.Note,
		reqDto.Image,
		reqDto.LanguageConfig,
		reqDto.FolderID,
	)

	clusterID, err := p.ps.Commands.CreateCluster.Handle(ctx, command)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	return httpPkg.SuccessCtxResponse(c, http.StatusCreated, "Cluster created", *clusterID)
}

// UpdateCluster
// @Tags clusters
// @Summary Update Cluster
// @Description Update existing Cluster
// @Accept json
// @Produce json
// @Param id path string true "Cluster ID"
// @Param Cluster body dto.UpdateClusterReqDto true "update Cluster"
// @Success 200 {string} id ""
// @Router /clusters/{id} [post]
func (p *clusterHandlers) UpdateCluster(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var reqDto requests.UpdateClusterReqDto
	if err := c.BodyParser(&reqDto); err != nil {
		p.log.Errorf("(Bind) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}
	err := p.val.DataValidation(reqDto)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	clusterID, err := primitive.ObjectIDFromHex(reqDto.ID)
	if err != nil {
		p.log.Errorf("(Handlers.GetByID)(uuid.FromString) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	command := clusterCommands.NewUpdateClusterCommand(
		clusterID.Hex(),
		reqDto.ClusterName,
		reqDto.Title,
		reqDto.Note,
		reqDto.Image,
		reqDto.LanguageConfig,
		reqDto.FolderID,
	)
	err = p.ps.Commands.UpdateCluster.Handle(ctx, command)

	if err != nil {
		p.log.Errorf("(Update.Handle) id: {%s}, err: {%v}", clusterID.Hex(), err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	p.log.Infof("(Cluster updated) id: {%s}", clusterID.Hex())
	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Cluster updated", clusterID.Hex())
}

// GetAllCluster
// @Tags clusters
// @Summary Get all clusters
// @Description Get all clusters
// @Accept json
// @Produce json
// @Success 200 {object} responses.GetAllClusterResponseDto
// @Router /clusters/ [get]
func (p *clusterHandlers) GetAllCluster(c *fiber.Ctx) error {
	ctx := c.Context()

	pq := utils.NewPaginationQuery(0, 0)

	response, err := p.ps.Queries.GetAllCluster.Handle(ctx, pq)
	if err != nil {
		p.log.Errorf("(Create.Handle) Error fetching clusters: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	p.log.Infof("(Hanlders.GetAll) result: {%+v}", response)

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Cluster found", response)
}

// GetClusterByID
// @Tags clusters
// @Summary Get Cluster
// @Description Get Cluster by id
// @Accept json
// @Produce json
// @Param id path string true "Cluster ID"
// @Success 200 {object} dto.ClusterResponseDto
// @Router /clusters/{id} [get]
func (p *clusterHandlers) GetClusterByID(c *fiber.Ctx) error {
	ctx := c.Context()
	param := c.Params(constants.ID)
	p.log.Infof("(Handlers.GetByID) id: {%s}", param)

	clusterID, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		p.log.Errorf("(Handlers.GetByID)(uuid.FromString) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	clusterQuery := clusterQueries.NewGetClusterByIDQuery(clusterID.Hex())
	err = p.val.DataValidation(clusterQuery)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	cluster, err := p.ps.Queries.GetClusterByID.Handle(ctx, clusterQuery)
	if err != nil {
		p.log.Errorf("(Handlers.GetByID)(Handle) id: {%s}, err: {%v}", clusterID.String(), err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	p.log.Infof("(Handlers.GetByID) clusterID: {%s}", clusterID.String())

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Cluster found", cluster)
}

// SearchCluster
// @Tags clusters
// @Summary Search clusters
// @Description Full text search by title and description
// @Accept json
// @Produce json
// @Param search queries string false "search text"
// @Param page queries string false "page number"
// @Param size queries string false "number of elements"
// @Success 200 {object} dto.ClusterSearchResponseDto
// @Router /clusters/search [get]
func (p *clusterHandlers) SearchCluster(c *fiber.Ctx) error {
	ctx := c.Context()
	pq := utils.NewPaginationFromQueryParams(c.Query(constants.Size), c.Query(constants.Page))

	var reqDto requests.SearchClusterFilterReqDto
	if err := c.QueryParser(&reqDto); err != nil {
		p.log.Errorf("(Bind) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}
	err := p.val.DataValidation(reqDto)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	clusterQuery := clusterQueries.NewSearchClustersQuery(
		reqDto.Keyword,
		pq,
	)

	searchRes, err := p.ps.Queries.SearchClusters.Handle(ctx, clusterQuery)
	if err != nil {
		p.log.Errorf("(Handlers.Search)(Handle) query: {%v}, err: {%v}", reqDto, err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	p.log.Infof("(Hanlders.Search) result: {%+v}", searchRes)

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Cluster found", searchRes)
}

// DeleteCluster
// @Tags clusters
// @Summary Delete Cluster
// @Description Delete Cluster by id
// @Accept json
// @Produce json
// @Param id path string true "Cluster ID"
// @Success 200 {object} dto.ClusterResponseDto
// @Router /clusters/{id} [post]
func (p *clusterHandlers) DeleteCluster(c *fiber.Ctx) error {
	ctx := c.Context()
	param := c.Params(constants.ID)
	p.log.Infof("(Handlers.Delete) id: {%s}", param)

	clusterID, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		p.log.Errorf("(Handlers.Delete)(uuid.FromString) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	clusterCommand := clusterCommands.NewDeleteClusterCommand(clusterID.Hex())
	err = p.val.DataValidation(clusterCommand)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	err = p.ps.Commands.DeleteCluster.Handle(ctx, clusterCommand)
	if err != nil {
		p.log.Errorf("(Handlers.Delete)(Handle) id: {%s}, err: {%v}", clusterID.String(), err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	p.log.Infof("(Handlers.Delete) clusterID: {%s}", clusterID.String())

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Cluster deleted", clusterID)
}

func (p *clusterHandlers) GetClusterComponents(c *fiber.Ctx) error {
	res := make([]responses.KeyValueResponseDto, 0, len(constants2.Components))
	for k, v := range constants2.Components {
		res = append(res, responses.KeyValueResponseDto{
			Key:   k,
			Value: v.String(),
		})
	}

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Cluster components found", res)
}

func (p *clusterHandlers) GetClusterLanguages(c *fiber.Ctx) error {
	res := make([]responses.KeyValueResponseDto, 0, len(constants2.GalleryLanguages))
	for k, v := range constants2.GalleryLanguages {
		res = append(res, responses.KeyValueResponseDto{
			Key:   k,
			Value: v.String(),
		})
	}

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Cluster languages found", res)
}
