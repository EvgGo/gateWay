package profile

import (
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func ListSkillsHandler(log *slog.Logger, c authv1.SkillsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "ListSkillsHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		reqLog.Info("Получен HTTP-запрос на получение списка навыков")

		q := strings.TrimSpace(r.URL.Query().Get("query"))
		pageToken := strings.TrimSpace(r.URL.Query().Get("page_token"))

		pageSize := int32(10) // дефолт для autocomplete
		if v := strings.TrimSpace(r.URL.Query().Get("page_size")); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil || n <= 0 {
				reqLog.Warn(
					"Некорректный параметр page_size: ожидалось положительное целое число",
					"page_size_raw", v,
					"err", err,
				)
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "page_size must be a positive integer")
				return
			}

			if n > 20 {
				reqLog.Info(
					"Параметр page_size превышает допустимый максимум, значение будет ограничено",
					"page_size_requested", n,
					"page_size_max", 20,
				)
				n = 20
			}
			pageSize = int32(n)
		}

		reqLog.Info(
			"Подготовлены параметры для запроса списка навыков",
			"query", q,
			"page_token", pageToken,
			"page_size", pageSize,
		)

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info(
			"Отправляем gRPC-запрос на получение списка навыков",
			"grpc_method", "ListSkills",
		)

		resp, err := c.ListSkills(ctx, &authv1.ListSkillsRequest{
			Query:     q,
			PageSize:  pageSize,
			PageToken: pageToken,
		})
		if err != nil {
			reqLog.Warn(
				"Ошибка при получении списка навыков из gRPC-сервиса",
				"grpc_method", "ListSkills",
				"query", q,
				"page_token", pageToken,
				"page_size", pageSize,
				"err", err,
			)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info(
			"Список навыков успешно получен и будет возвращен клиенту",
			"grpc_method", "ListSkills",
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func GetSkillHandler(log *slog.Logger, c authv1.SkillsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "GetSkillHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
		)

		reqLog.Info("Получен HTTP-запрос на получение навыка по идентификатору")

		skillID := strings.TrimSpace(chi.URLParam(r, "skill_id"))
		if skillID == "" {
			reqLog.Warn("Не передан обязательный path-параметр skill_id")
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "skill_id is required")
			return
		}

		reqLog.Info(
			"Подготовлены параметры для запроса навыка",
			"skill_id", skillID,
		)

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info(
			"Отправляем gRPC-запрос на получение навыка",
			"grpc_method", "GetSkill",
			"skill_id", skillID,
		)

		resp, err := c.GetSkill(ctx, &authv1.GetSkillRequest{SkillId: skillID})
		if err != nil {
			reqLog.Warn(
				"Ошибка при получении навыка из gRPC-сервиса",
				"grpc_method", "GetSkill",
				"skill_id", skillID,
				"err", err,
			)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info(
			"Навык успешно получен и будет возвращен клиенту",
			"grpc_method", "GetSkill",
			"skill_id", skillID,
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func CreateSkillHandler(log *slog.Logger, c authv1.SkillsClient) http.HandlerFunc {

	type createSkillJSON struct {
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "CreateSkillHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
		)

		reqLog.Info("Получен HTTP-запрос на создание нового навыка")

		var in createSkillJSON
		if err := utils.DecodeJSON(w, r, &in); err != nil {
			reqLog.Warn(
				"Не удалось декодировать JSON-тело запроса при создании навыка",
				"err", err,
			)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid JSON body")
			return
		}

		name := strings.TrimSpace(in.Name)
		if name == "" {
			reqLog.Warn("В теле запроса не заполнено обязательное поле name")
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "name is required")
			return
		}

		reqLog.Info(
			"Подготовлены данные для создания навыка",
			"name", name,
		)

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info(
			"Отправляем gRPC-запрос на создание навыка",
			"grpc_method", "CreateSkill",
			"name", name,
		)

		resp, err := c.CreateSkill(ctx, &authv1.CreateSkillRequest{Name: name})
		if err != nil {
			reqLog.Warn(
				"Ошибка при создании навыка в gRPC-сервисе",
				"grpc_method", "CreateSkill",
				"name", name,
				"err", err,
			)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info(
			"Навык успешно создан и будет возвращен клиенту",
			"grpc_method", "CreateSkill",
			"name", name,
		)

		helpers.WriteProtoJSON(w, http.StatusCreated, resp)
	}
}

func DeleteSkillHandler(log *slog.Logger, c authv1.SkillsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "DeleteSkillHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
		)

		reqLog.Info("Получен HTTP-запрос на удаление навыка")

		skillID := strings.TrimSpace(chi.URLParam(r, "skill_id"))
		if skillID == "" {
			reqLog.Warn("Не передан обязательный path-параметр skill_id для удаления навыка")
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "skill_id is required")
			return
		}

		reqLog.Info(
			"Подготовлены параметры для удаления навыка",
			"skill_id", skillID,
		)

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info(
			"Отправляем gRPC-запрос на удаление навыка",
			"grpc_method", "DeleteSkill",
			"skill_id", skillID,
		)

		_, err := c.DeleteSkill(ctx, &authv1.DeleteSkillRequest{SkillId: skillID})
		if err != nil {
			reqLog.Warn(
				"Ошибка при удалении навыка в gRPC-сервисе",
				"grpc_method", "DeleteSkill",
				"skill_id", skillID,
				"err", err,
			)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info(
			"Навык успешно удален",
			"grpc_method", "DeleteSkill",
			"skill_id", skillID,
		)

		w.WriteHeader(http.StatusNoContent)
	}
}

type GetSkillsByIdsJSON struct {
	SkillIDs []string `json:"skill_ids"`
}

func GetSkillsByIdsHandler(log *slog.Logger, c authv1.SkillsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "GetSkillsByIdsHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
		)

		reqLog.Info("Получен HTTP-запрос на получение навыков по идентификаторам")

		var in GetSkillsByIdsJSON
		if err := utils.DecodeJSON(w, r, &in); err != nil {
			reqLog.Warn(
				"Не удалось декодировать JSON-тело запроса",
				"err", err,
			)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid JSON body")
			return
		}

		uniqueIDs := make([]string, 0, len(in.SkillIDs))
		seen := make(map[string]struct{}, len(in.SkillIDs))

		for _, rawID := range in.SkillIDs {
			skillID := strings.TrimSpace(rawID)
			if skillID == "" {
				continue
			}
			if _, exists := seen[skillID]; exists {
				continue
			}
			seen[skillID] = struct{}{}
			uniqueIDs = append(uniqueIDs, skillID)
		}

		if len(uniqueIDs) == 0 {
			reqLog.Warn("Не передан ни один корректный skill_id")
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "skill_ids is required")
			return
		}

		reqLog.Info(
			"Подготовлены параметры для batch-запроса навыков",
			"skill_ids_count", len(uniqueIDs),
			"skill_ids", uniqueIDs,
		)

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info(
			"Отправляем gRPC-запрос на получение навыков по ids",
			"grpc_method", "GetSkillsByIds",
			"skill_ids_count", len(uniqueIDs),
		)

		resp, err := c.GetSkillsByIds(ctx, &authv1.GetSkillsByIdsRequest{
			SkillIds: uniqueIDs,
		})
		if err != nil {
			reqLog.Warn(
				"Ошибка при получении навыков из gRPC-сервиса",
				"grpc_method", "GetSkillsByIds",
				"skill_ids_count", len(uniqueIDs),
				"err", err,
			)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info(
			"Навыки успешно получены и будут возвращены клиенту",
			"grpc_method", "GetSkillsByIds",
			"requested_skill_ids_count", len(uniqueIDs),
			"returned_skills_count", len(resp.GetSkills()),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
