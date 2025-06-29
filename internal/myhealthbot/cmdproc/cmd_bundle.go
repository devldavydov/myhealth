package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/devldavydov/myhealth/internal/common/html"
	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (r *CmdProcessor) bundleSetCommand(userID int64, key string, bndlParts []string) []CmdResponse {
	bndlData := make(map[string]float64)

	for _, bndlPart := range bndlParts {
		if strings.Contains(bndlPart, ":") {
			// Add dependant food
			parts := strings.Split(bndlPart, ":")
			if len(parts) > 2 {
				return NewSingleCmdResponse(MsgErrInvalidCommand)
			}

			weight, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return NewSingleCmdResponse(MsgErrInvalidCommand)
			}

			bndlData[parts[0]] = weight
		} else {
			// Add dependant bundle key.
			bndlData[bndlPart] = 0
		}
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetBundle(ctx, userID, &storage.Bundle{Key: key, Data: bndlData}, true); err != nil {
		if errors.Is(err, storage.ErrBundleInvalid) {
			return NewSingleCmdResponse(MsgErrInvalidCommand)
		}
		if errors.Is(err, storage.ErrBundleDepBundleNotFound) {
			return NewSingleCmdResponse(MsgErrBundleDepBundleNotFound)
		}
		if errors.Is(err, storage.ErrBundleDepRecursive) {
			return NewSingleCmdResponse(MsgErrBundleDepBundleRecursive)
		}
		if errors.Is(err, storage.ErrBundleDepFoodNotFound) {
			return NewSingleCmdResponse(MsgErrBundleDepFoodNotFound)
		}

		r.logger.Error(
			"bundle set command DB error",
			zap.Strings("cmdParts", bndlParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) bundleSetTemplateCommand(userID int64, key string) []CmdResponse {
	// Get in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	bndl, err := r.stg.GetBundle(ctx, userID, key)
	if err != nil {
		if errors.Is(err, storage.ErrBundleNotFound) {
			return NewSingleCmdResponse(MsgErrBundleNotFound)
		}

		r.logger.Error(
			"bundle set template command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("b,set,%s,", bndl.Key))

	items := make([]string, 0, len(bndl.Data))
	for k, v := range bndl.Data {
		if v > 0 {
			items = append(items, fmt.Sprintf("%s:%1.f", k, v))
		} else {
			items = append(items, k)
		}
	}
	sb.WriteString(strings.Join(items, "/"))

	return NewSingleCmdResponse(sb.String())
}

func (r *CmdProcessor) bundleListCommand(userID int64) []CmdResponse {
	// Get in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	lst, err := r.stg.GetBundleList(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(MsgErrEmptyResult)
		}

		r.logger.Error(
			"bundle list command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	// Build html
	htmlBuilder := html.NewBuilder("Список бандлов")

	// Table
	tbl := html.NewTable([]string{
		"Ключ бандла", "Еда/Ключ дочернего бандла", "Вес еды, г.",
	})

	for _, bndl := range lst {
		i := 0
		for k, v := range bndl.Data {
			tr := html.NewTr(nil)
			if i == 0 {
				tr.AddTd(html.NewTd(html.NewS(bndl.Key), html.Attrs{"rowspan": strconv.Itoa(len(bndl.Data))}))
			}
			if v > 0 {
				tr.AddTd(html.NewTd(html.NewS(k), nil))
				tr.AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.1f", v)), nil))
			} else {
				tr.AddTd(html.NewTd(html.NewI(k, nil), nil))
				tr.AddTd(html.NewTd(html.NewS(""), nil))
			}
			i++
			tbl.AddRow(tr)
		}
	}

	// Doc
	htmlBuilder.Add(
		html.NewContainer().Add(
			html.NewH(
				"Список бандлов",
				5,
				html.Attrs{"align": "center"},
			),
			tbl))

	// Response
	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: "bundles.html",
	})
}

func (r *CmdProcessor) bundleDelCommand(userID int64, key string) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteBundle(ctx, userID, key); err != nil {
		if errors.Is(err, storage.ErrBundleIsUsed) {
			return NewSingleCmdResponse(MsgErrBundleIsUsed)
		}

		r.logger.Error(
			"bundle del command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}
