package bot

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (k *RealBot) saveVideo(chatID int64, fileID string) {
	// Получение файлов?
	fmt.Println("info")

	name, file, err := tgbotapi.NewVideo(chatID, tgbotapi.FileID(fileID)).File.UploadData()
	if err != nil {
		k.bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Файл %s не был загружен: %s", name, err.Error())))
		return
	}
	fmt.Println("read")

	b, err := io.ReadAll(file)
	if err != nil {
		k.bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Файл %s не может быть прочитан: %s", name, err.Error())))
		return
	}
	fmt.Println("write")

	if err := os.WriteFile(k.downloadPath+"/"+name, b, os.FileMode(os.O_CREATE)); err != nil {
		k.bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Файл %s не был сохранен: %s", name, err.Error())))
		return
	}
}

func (k *RealBot) loadVideo(chatID int64, fileName string) error {
	local_video, err := os.ReadFile(fmt.Sprintf("%s/%s", k.downloadPath, fileName))
	if err != nil {
		return fmt.Errorf("ReadFile: %w", err)
	}

	videoFileBytes := tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: local_video,
	}

	k.bot.Send(tgbotapi.NewVideo(chatID, videoFileBytes))

	return nil
}

func (k *RealBot) removeVideo(chatID int64, fileName string) {
	if err := os.Remove(fmt.Sprintf("%s/%s", k.downloadPath, fileName)); err != nil {
		k.bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Не удается удалить файл %s с сервера: %s", fileName, err.Error())))
	}
	k.bot.Send(tgbotapi.NewMessage(chatID, "Видео удалено"))
}

func (k *RealBot) toAdmins(msg string) {
	for v := range k.admins {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			fmt.Println("ParseintError:", err)
		}
		k.bot.Send(tgbotapi.NewMessage(id, msg))
	}
}

func (k *RealBot) getAllowedFiles(mode MODE) ([]string, error) {
	files, err := os.ReadDir(k.downloadPath)
	if err != nil {
		return nil, fmt.Errorf("caanot read directory: %w", err)
	}

	res := make([]string, 0, len(files))
	var info fs.FileInfo
	for i := range files {
		info, err = files[i].Info()
		if (err == nil && (info.Size()/1024/1024) < 2048) || mode == RemoveMode {
			res = append(res, files[i].Name())
		}
	}

	return res, nil
}
