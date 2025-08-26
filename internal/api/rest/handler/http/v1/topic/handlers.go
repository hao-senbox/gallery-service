package topic

import (
	"gallery-service/config"
	"gallery-service/internal/api/rest/validator"
	topicCommands "gallery-service/internal/application/commands/v1/topic"
	requests "gallery-service/internal/application/dto/requests/topic"
	"gallery-service/internal/application/dto/responses"
	topicQueries "gallery-service/internal/application/queries/topic"
	"gallery-service/internal/domain/service"
	constants2 "gallery-service/internal/pkg/constants"
	"gallery-service/pkg/constants"
	httpPkg "gallery-service/pkg/http"
	"gallery-service/pkg/utils"
	"gallery-service/pkg/zap"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type topicHandlers struct {
	log         zap.Logger
	cfg         *config.Config
	ps          *service.TopicService
	val         *validator.Wrapper
	mongoClient *mongo.Client
}

func NewTopicHandlers(
	log zap.Logger,
	cfg *config.Config,
	mongoClient *mongo.Client,
) *topicHandlers {
	return &topicHandlers{
		log:         log,
		cfg:         cfg,
		val:         validator.NewValidator(log, cfg),
		mongoClient: mongoClient,
	}
}

// CreateTopic
// @Tags Topics
// @Summary Create Topic
// @Description Create new Topic
// @Param Topic body dto.CreateTopicReqDto true "create Topic"
// @Accept json
// @Produce json
// @Success 201 {string} id ""
// @Router /topics [post]
func (p *topicHandlers) CreateTopic(c *fiber.Ctx) error {
	ctx := c.Context()
	var reqDto requests.CreateTopicReqDto
	if err := c.BodyParser(&reqDto); err != nil {
		p.log.Errorf("(Bind) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}
	err := p.val.DataValidation(reqDto)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	command := topicCommands.NewCreateTopicCommand(
		reqDto.TopicName,
		reqDto.IsPublished,
		reqDto.LanguageConfig,
	)

	topicID, err := p.ps.Commands.CreateTopic.Handle(ctx, command)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	return httpPkg.SuccessCtxResponse(c, http.StatusCreated, "Topic created", *topicID)
}

// UpdateTopic
// @Tags Topics
// @Summary Update Topic
// @Description Update existing Topic
// @Accept json
// @Produce json
// @Param id path string true "Topic ID"
// @Param Topic body dto.UpdateTopicReqDto true "update Topic"
// @Success 200 {string} id ""
// @Router /topics/{id} [post]
func (p *topicHandlers) UpdateTopic(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var reqDto requests.UpdateTopicReqDto
	if err := c.BodyParser(&reqDto); err != nil {
		p.log.Errorf("(Bind) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}
	err := p.val.DataValidation(reqDto)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	topicID, err := primitive.ObjectIDFromHex(reqDto.ID)
	if err != nil {
		p.log.Errorf("(Handlers.GetByID)(uuid.FromString) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	command := topicCommands.NewUpdateTopicCommand(
		topicID.Hex(),
		reqDto.FileName,
		reqDto.IsPublished,
		reqDto.LanguageConfig,
	)
	err = p.ps.Commands.UpdateTopic.Handle(ctx, command)

	if err != nil {
		p.log.Errorf("(Update.Handle) id: {%s}, err: {%v}", topicID.Hex(), err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	p.log.Infof("(Topic updated) id: {%s}", topicID.Hex())
	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Topic updated", topicID.Hex())
}

// GetAllTopic
// @Tags topics
// @Summary Get all Topics
// @Description Get all Topics
// @Accept json
// @Produce json
// @Success 200 {object} responses.GetAllTopicResponseDto
// @Router /topics/ [get]
func (p *topicHandlers) GetAllTopic(c *fiber.Ctx) error {
	ctx := c.Context()

	pq := utils.NewPaginationQuery(0, 0)

	response, err := p.ps.Queries.GetAllTopic.Handle(ctx, pq)
	if err != nil {
		p.log.Errorf("(Create.Handle) Error fetching topics: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Topic found", response)
}

// GetTopicByID
// @Tags Topics
// @Summary Get Topic
// @Description Get Topic by id
// @Accept json
// @Produce json
// @Param id path string true "Topic ID"
// @Success 200 {object} dto.TopicResponseDto
// @Router /topics/{id} [get]
func (p *topicHandlers) GetTopicByID(c *fiber.Ctx) error {
	ctx := c.Context()
	param := c.Params(constants.ID)
	p.log.Infof("(Handlers.GetByID) id: {%s}", param)

	topicID, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		p.log.Errorf("(Handlers.GetByID)(uuid.FromString) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	topicQuery := topicQueries.NewGetTopicByIDQuery(topicID.Hex())
	err = p.val.DataValidation(topicQuery)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	topic, err := p.ps.Queries.GetTopicByID.Handle(ctx, topicQuery)
	if err != nil {
		p.log.Errorf("(Handlers.GetByID)(Handle) id: {%s}, err: {%v}", topicID.String(), err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Topic found", topic)
}

// SearchTopic
// @Tags Topics
// @Summary Search Topics
// @Description Full text search by title and description
// @Accept json
// @Produce json
// @Param search queries string false "search text"
// @Param page queries string false "page number"
// @Param size queries string false "number of elements"
// @Success 200 {object} dto.TopicSearchResponseDto
// @Router /topics/search [get]
func (p *topicHandlers) SearchTopic(c *fiber.Ctx) error {
	ctx := c.Context()
	pq := utils.NewPaginationFromQueryParams(c.Query(constants.Size), c.Query(constants.Page))

	var reqDto requests.SearchTopicFilterReqDto
	if err := c.QueryParser(&reqDto); err != nil {
		p.log.Errorf("(Bind) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}
	err := p.val.DataValidation(reqDto)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	topicQuery := topicQueries.NewSearchTopicsQuery(
		reqDto.Keyword,
		pq,
	)

	searchRes, err := p.ps.Queries.SearchTopics.Handle(ctx, topicQuery)
	if err != nil {
		p.log.Errorf("(Handlers.Search)(Handle) query: {%v}, err: {%v}", reqDto, err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Topic found", searchRes)
}

// DeleteTopic
// @Tags Topics
// @Summary Delete Topic
// @Description Delete Topic by id
// @Accept json
// @Produce json
// @Param id path string true "Topic ID"
// @Success 200 {object} dto.TopicResponseDto
// @Router /topics/{id} [post]
func (p *topicHandlers) DeleteTopic(c *fiber.Ctx) error {
	ctx := c.Context()
	param := c.Params(constants.ID)
	p.log.Infof("(Handlers.Delete) id: {%s}", param)

	topicID, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		p.log.Errorf("(Handlers.Delete)(uuid.FromString) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	topicCommand := topicCommands.NewDeleteTopicCommand(topicID.Hex())
	err = p.val.DataValidation(topicCommand)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	err = p.ps.Commands.DeleteTopic.Handle(ctx, topicCommand)
	if err != nil {
		p.log.Errorf("(Handlers.Delete)(Handle) id: {%s}, err: {%v}", topicID.String(), err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	p.log.Infof("(Handlers.Delete) topicID: {%s}", topicID.String())

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Topic deleted", topicID)
}

func (p *topicHandlers) GetTopicComponents(c *fiber.Ctx) error {
	res := make([]responses.KeyValueResponseDto, 0, len(constants2.Components))
	for k, v := range constants2.Components {
		res = append(res, responses.KeyValueResponseDto{
			Key:   k,
			Value: v.String(),
		})
	}

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Topic components found", res)
}

func (p *topicHandlers) GetTopicLanguages(c *fiber.Ctx) error {
	res := make([]responses.KeyValueResponseDto, 0, len(constants2.GalleryLanguages))
	for k, v := range constants2.GalleryLanguages {
		res = append(res, responses.KeyValueResponseDto{
			Key:   k,
			Value: v.String(),
		})
	}

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Topic languages found", res)
}

func (p *topicHandlers) GetAllTopic4App(c *fiber.Ctx) error {
	ctx := c.Context()
	response, err := p.ps.Queries.GetAllTopic.Handle4App(ctx)
	if err != nil {
		p.log.Errorf("(Create.Handle) Error fetching topics: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Topic found", response)
}
