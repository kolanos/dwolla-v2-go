package dwolla

import (
	"context"
	"errors"
	"fmt"
)

// KBAService is the kba service interface
//
// see: https://developers.dwolla.com/api-reference/kba
type KBAService interface {
	Retrieve(context.Context, string) (*KBA, error)
}

// KBAServiceOp is an implementation of the kba service interface
type KBAServiceOp struct {
	client *Client
}

// KBA is a knowledge based authentication resource
type KBA struct {
	Resource
	ID        string        `json:"id"`
	Questions []KBAQuestion `json:"questions"`
}

// KBAQuestion is a knowledge based authentication question
type KBAQuestion struct {
	ID      string      `json:"id"`
	Text    string      `json:"text"`
	Answers []KBAAnswer `json:"answers"`
}

// KBAAnswer is a knowledge based authentication answer
type KBAAnswer struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// KBARequest is a knowledge based authentication verification request
type KBARequest struct {
	Answers []KBAQuestionAnswer `json:"answers"`
}

// KBAQuestionAnswer is a knowledge based authentication question and answer
type KBAQuestionAnswer struct {
	QuestionID string `json:"questionId"`
	AnswerID   string `json:"answerId"`
}

// Retrieve retrieves a knowledge based authentication session
//
// see: https://docs.dwolla.com/#retrieve-kba-questions
func (k *KBAServiceOp) Retrieve(ctx context.Context, id string) (*KBA, error) {
	var kba KBA

	if err := k.client.Get(ctx, fmt.Sprintf("kba/%s", id), nil, nil, &kba); err != nil {
		return nil, err
	}

	kba.client = k.client

	return &kba, nil
}

// Verify attempts a knowledge based authentication verification
//
// see: https://docs.dwolla.com/#verify-kba-questions
func (k *KBA) Verify(ctx context.Context, body *KBARequest) error {
	if _, ok := k.Links["self"]; !ok {
		return errors.New("No self resource link")
	}

	return k.client.Post(ctx, k.Links["self"].Href, body, nil, nil)
}
