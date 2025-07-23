package folder

import (
	"gallery-service/config"
	"gallery-service/internal/api/rest/validator"
	folderCommands "gallery-service/internal/application/commands/v1/folder"
	requests "gallery-service/internal/application/dto/requests/folder"
	folderQueries "gallery-service/internal/application/queries/folder"
	"gallery-service/internal/domain/service"
	"gallery-service/pkg/constants"
	httpPkg "gallery-service/pkg/http"
	"gallery-service/pkg/utils"
	"gallery-service/pkg/zap"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type folderHandlers struct {
	log         zap.Logger
	cfg         *config.Config
	ps          *service.FolderService
	val         *validator.Wrapper
	mongoClient *mongo.Client
}

func NewFolderHandlers(
	log zap.Logger,
	cfg *config.Config,
	mongoClient *mongo.Client,
) *folderHandlers {
	return &folderHandlers{
		log:         log,
		cfg:         cfg,
		val:         validator.NewValidator(log, cfg),
		mongoClient: mongoClient,
	}
}

// CreateFolder
// @Tags folders
// @Summary Create Folder
// @Description Create new Folder
// @Param Folder body dto.CreateFolderReqDto true "create Folder"
// @Accept json
// @Produce json
// @Success 201 {string} id ""
// @Router /folders [post]
func (p *folderHandlers) CreateFolder(c *fiber.Ctx) error {
	ctx := c.Context()
	var reqDto requests.CreateFolderReqDto
	if err := c.BodyParser(&reqDto); err != nil {
		p.log.Errorf("(Bind) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}
	err := p.val.DataValidation(reqDto)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	command := folderCommands.NewCreateFolderCommand(
		reqDto.FolderName,
		reqDto.FolderThumbnailKey,
		reqDto.FolderThumbnailURL,
		reqDto.ParentID,
	)

	folderID, err := p.ps.Commands.CreateFolder.Handle(ctx, command)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	return httpPkg.SuccessCtxResponse(c, http.StatusCreated, "Folder created", *folderID)
}

// UpdateFolder
// @Tags folders
// @Summary Update Folder
// @Description Update existing Folder
// @Accept json
// @Produce json
// @Param id path string true "Folder ID"
// @Param Folder body dto.UpdateFolderReqDto true "update Folder"
// @Success 200 {string} id ""
// @Router /folders/{id} [post]
func (p *folderHandlers) UpdateFolder(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var reqDto requests.UpdateFolderReqDto
	if err := c.BodyParser(&reqDto); err != nil {
		p.log.Errorf("(Bind) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}
	err := p.val.DataValidation(reqDto)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	folderID, err := primitive.ObjectIDFromHex(reqDto.ID)
	if err != nil {
		p.log.Errorf("(Handlers.GetByID)(uuid.FromString) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	command := folderCommands.NewUpdateFolderCommand(
		folderID.Hex(),
		reqDto.FolderName,
		reqDto.FolderThumbnailKey,
		reqDto.FolderThumbnailURL,
		reqDto.ParentID,
	)
	err = p.ps.Commands.UpdateFolder.Handle(ctx, command)

	if err != nil {
		p.log.Errorf("(Update.Handle) id: {%s}, err: {%v}", folderID.Hex(), err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	p.log.Infof("(Folder updated) id: {%s}", folderID.Hex())
	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Folder updated", folderID.Hex())
}

// GetAllFolder
// @Tags folders
// @Summary Get all folders
// @Description Get all folders
// @Accept json
// @Produce json
// @Success 200 {object} responses.GetAllFolderResponseDto
// @Router /folders/ [get]
func (p *folderHandlers) GetAllFolder(c *fiber.Ctx) error {
	ctx := c.Context()

	pq := utils.NewPaginationQuery(0, 0)

	response, err := p.ps.Queries.GetAllFolder.Handle(ctx, pq)
	if err != nil {
		p.log.Errorf("(Create.Handle) Error fetching folders: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	p.log.Infof("(Hanlders.GetAll) result: {%+v}", response)

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Folder found", response)
}

// GetFolderByID
// @Tags folders
// @Summary Get Folder
// @Description Get Folder by id
// @Accept json
// @Produce json
// @Param id path string true "Folder ID"
// @Success 200 {object} dto.FolderResponseDto
// @Router /folders/{id} [get]
func (p *folderHandlers) GetFolderByID(c *fiber.Ctx) error {
	ctx := c.Context()
	param := c.Params(constants.ID)
	p.log.Infof("(Handlers.GetByID) id: {%s}", param)

	folderID, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		p.log.Errorf("(Handlers.GetByID)(uuid.FromString) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	folderQuery := folderQueries.NewGetFolderByIDQuery(folderID.Hex())
	err = p.val.DataValidation(folderQuery)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	folder, err := p.ps.Queries.GetFolderByID.Handle(ctx, folderQuery)
	if err != nil {
		p.log.Errorf("(Handlers.GetByID)(Handle) id: {%s}, err: {%v}", folderID.String(), err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	p.log.Infof("(Handlers.GetByID) folderID: {%s}", folderID.String())

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Folder found", folder)
}

// SearchFolder
// @Tags folders
// @Summary Search folders
// @Description Full text search by title and description
// @Accept json
// @Produce json
// @Param search queries string false "search text"
// @Param page queries string false "page number"
// @Param size queries string false "number of elements"
// @Success 200 {object} dto.FolderSearchResponseDto
// @Router /folders/search [get]
func (p *folderHandlers) SearchFolder(c *fiber.Ctx) error {
	ctx := c.Context()
	pq := utils.NewPaginationFromQueryParams(c.Query(constants.Size), c.Query(constants.Page))

	var reqDto requests.SearchFolderFilterReqDto
	if err := c.QueryParser(&reqDto); err != nil {
		p.log.Errorf("(Bind) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}
	err := p.val.DataValidation(reqDto)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	folderQuery := folderQueries.NewSearchFoldersQuery(
		reqDto.Keyword,
		pq,
	)

	searchRes, err := p.ps.Queries.SearchFolders.Handle(ctx, folderQuery)
	if err != nil {
		p.log.Errorf("(Handlers.Search)(Handle) query: {%v}, err: {%v}", reqDto, err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	p.log.Infof("(Hanlders.Search) result: {%+v}", searchRes)

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Folder found", searchRes)
}

// DeleteFolder
// @Tags folders
// @Summary Delete Folder
// @Description Delete Folder by id
// @Accept json
// @Produce json
// @Param id path string true "Folder ID"
// @Success 200 {object} dto.FolderResponseDto
// @Router /folders/{id} [post]
func (p *folderHandlers) DeleteFolder(c *fiber.Ctx) error {
	ctx := c.Context()
	param := c.Params(constants.ID)
	p.log.Infof("(Handlers.Delete) id: {%s}", param)

	folderID, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		p.log.Errorf("(Handlers.Delete)(uuid.FromString) err: {%v}", err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	folderCommand := folderCommands.NewDeleteFolderCommand(folderID.Hex())
	err = p.val.DataValidation(folderCommand)
	if err != nil {
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	err = p.ps.Commands.DeleteFolder.Handle(ctx, folderCommand)
	if err != nil {
		p.log.Errorf("(Handlers.Delete)(Handle) id: {%s}, err: {%v}", folderID.String(), err)
		return httpPkg.ErrorCtxResponse(c, err, p.cfg.App.API.Rest.Setting.DebugErrorsResponse)
	}

	p.log.Infof("(Handlers.Delete) folderID: {%s}", folderID.String())

	return httpPkg.SuccessCtxResponse(c, http.StatusOK, "Folder deleted", folderID)
}
