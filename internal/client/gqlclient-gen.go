// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package client

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type ClientID string

const (
	ClientIDServarr ClientID = "Servarr"
)

// ServicesDownloadClientClientMutation includes the requested fields of the GraphQL type ClientMutation.
type ServicesDownloadClientClientMutation struct {
	Download string `json:"download"`
}

// GetDownload returns ServicesDownloadClientClientMutation.Download, and is useful for accessing the field via an interface.
func (v *ServicesDownloadClientClientMutation) GetDownload() string { return v.Download }

// ServicesDownloadResponse is returned by ServicesDownload on success.
type ServicesDownloadResponse struct {
	Client ServicesDownloadClientClientMutation `json:"client"`
}

// GetClient returns ServicesDownloadResponse.Client, and is useful for accessing the field via an interface.
func (v *ServicesDownloadResponse) GetClient() ServicesDownloadClientClientMutation { return v.Client }

// __ServicesDownloadInput is used internally by genqlient
type __ServicesDownloadInput struct {
	InfoHashes []protocol.ID `json:"infoHashes"`
	ClientID   ClientID      `json:"clientID"`
}

// GetInfoHashes returns __ServicesDownloadInput.InfoHashes, and is useful for accessing the field via an interface.
func (v *__ServicesDownloadInput) GetInfoHashes() []protocol.ID { return v.InfoHashes }

// GetClientID returns __ServicesDownloadInput.ClientID, and is useful for accessing the field via an interface.
func (v *__ServicesDownloadInput) GetClientID() ClientID { return v.ClientID }

// The query or mutation executed by ServicesDownload.
const ServicesDownload_Operation = `
mutation ServicesDownload ($infoHashes: [Hash20!], $clientID: ClientID) {
	client {
		download(infoHashes: $infoHashes, clientID: $clientID)
	}
}
`

func ServicesDownload(
	ctx_ context.Context,
	client_ graphql.Client,
	infoHashes []protocol.ID,
	clientID ClientID,
) (*ServicesDownloadResponse, error) {
	req_ := &graphql.Request{
		OpName: "ServicesDownload",
		Query:  ServicesDownload_Operation,
		Variables: &__ServicesDownloadInput{
			InfoHashes: infoHashes,
			ClientID:   clientID,
		},
	}
	var err_ error

	var data_ ServicesDownloadResponse
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return &data_, err_
}
