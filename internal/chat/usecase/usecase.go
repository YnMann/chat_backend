package usecase

import (
	"context"

	"github.com/YnMann/chat_backend/internal/chat"
	"github.com/YnMann/chat_backend/internal/models"
)

type ChatUseCase struct {
	chatRepo chat.ChatRepository
}

func NewChatUseCase(
	chatRepo chat.ChatRepository,
) *ChatUseCase {
	return &ChatUseCase{
		chatRepo: chatRepo,
	}
}

func (uc *ChatUseCase) SetUserOnlineStatus(ctx context.Context, uID string, isOnline bool) error {
	err := uc.chatRepo.SetUserOnlineStatus(ctx, uID, isOnline)
	if err != nil {
		return err
	}

	return nil
}

func (uc *ChatUseCase) CreateMsg(context.Context, *models.Messages) error {

	return nil
}

func (uc *ChatUseCase) GetMsg(ctx context.Context, sender string, sender_ip string, recipient string) (*models.Messages, error) {

	return nil, nil
}
