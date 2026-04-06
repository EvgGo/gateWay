package profile

import (
	"gateWay/internal/DTO/profile"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"log/slog"
	"net/http"
	"strings"
)

// UpdateMeHandler обрабатывает PATCH /users/me
//
// Алгоритм:
// 1) decode JSON в updateMeJSON
// 2) Собираем protobuf UpdateMeRequest, передавая optional поля как pointers
// 3) competence_levels разбираем только если ключ реально был в JSON
// 4) gRPC UpdateMe
// 5) отдаем обновленного User
func UpdateMeHandler(log *slog.Logger, c authv1.UserProfileClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "UpdateMeHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
		)

		reqLog.Info("Получен HTTP-запрос на обновление профиля текущего пользователя")

		var in profile.UpdateMeJSON
		if err := utils.DecodeJSON(w, r, &in); err != nil {
			reqLog.Warn(
				"Не удалось декодировать JSON-тело запроса при обновлении профиля",
				"err", err,
			)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid JSON body")
			return
		}

		req := &authv1.UpdateMeRequest{
			FirstName:             in.FirstName,
			LastName:              in.LastName,
			Phone:                 in.Phone,
			About:                 in.About,
			IsUserOpenSuggestions: in.IsUserOpenSuggestions,
			IsProfileHidden:       in.IsProfileHidden,
		}

		reqLog.Info(
			"Основные поля запроса на обновление профиля успешно разобраны",
			"first_name_len", len(strings.TrimSpace(utils.StringValue(in.FirstName))),
			"last_name_len", len(strings.TrimSpace(utils.StringValue(in.LastName))),
			"phone_present", strings.TrimSpace(utils.StringValue(in.Phone)) != "",
			"about_len", len(strings.TrimSpace(utils.StringValue(in.About))),
			"is_user_open_suggestions", in.IsUserOpenSuggestions,
			"is_profile_hidden", in.IsProfileHidden,
			"competence_levels_present", in.CompetenceLevels != nil,
			"skills_present", in.Skills != nil,
		)

		// competence_levels: обновляем только если ключ был в JSON
		// Если ключ был:
		// - {} => очистить карту (или заменить на пустую),
		// - {"go":"senior"} => заменить содержимое карты
		if in.CompetenceLevels != nil {
			reqLog.Info("В запросе присутствует поле competence_levels, начинаем разбор")

			m, err := parseCompetenceLevels(in.CompetenceLevels)
			if err != nil {
				reqLog.Warn(
					"Некорректное значение поля competence_levels: ожидался объект или null",
					"err", err,
				)
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "competence_levels must be an object or null")
				return
			}

			req.CompetenceLevels = m

			reqLog.Info(
				"Поле competence_levels успешно обработано",
				"competence_levels_count", len(m),
			)
		} else {
			reqLog.Info("Поле competence_levels отсутствует в запросе, текущее значение не будет изменено")
		}

		// skills:
		// - поле отсутствует -> не трогаем
		// - null -> очистить все
		// - [] -> очистить все
		// - {"ids":["1","2"]} -> заменить
		// - ["1","2"] -> заменить
		if in.Skills != nil {
			reqLog.Info("В запросе присутствует поле skills, начинаем разбор")

			skills, err := parseSkillsSelection(in.Skills)
			if err != nil {
				reqLog.Warn(
					"Некорректное значение поля skills",
					"err", err,
				)
				helpers.WriteAPIError(w, r, http.StatusBadRequest, err.Error())
				return
			}
			req.Skills = skills

			reqLog.Info(
				"Поле skills успешно обработано",
				"skills_count", len(skills.Ids),
			)
		} else {
			reqLog.Info("Поле skills отсутствует в запросе, текущий список навыков не будет изменен")
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info(
			"Отправляем gRPC-запрос на обновление профиля текущего пользователя",
			"grpc_method", "UpdateMe",
			"competence_levels_count", len(req.CompetenceLevels),
			"skills_count", len(req.Skills.Ids),
		)

		resp, err := c.UpdateMe(ctx, req)
		if err != nil {
			reqLog.Warn(
				"Ошибка при обновлении профиля в gRPC-сервисе",
				"grpc_method", "UpdateMe",
				"err", err,
			)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info(
			"Профиль текущего пользователя успешно обновлен и будет возвращен клиенту",
			"grpc_method", "UpdateMe",
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
