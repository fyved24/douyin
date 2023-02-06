// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    commentListResponse, err := UnmarshalCommentListResponse(bytes)
//    bytes, err = CommentListResponse.Marshal()
//    commentActionResponse, err := UnmarshalCommentActionResponse(bytes)
//    bytes, err = CommentActionResponse.Marshal()

package responses

import "encoding/json"

type CommentListResponse struct {
	CommonResponse
	CommentList []Comment `json:"comment_list"` // 评论列表
}

type CommentActionResponse struct {
	CommonResponse
	Comment Comment `json:"comment"` // 评论成功返回评论内容，不需要重新拉取整个列表
}

func UnmarshalCommentListResponse(data []byte) (CommentListResponse, error) {
	var r CommentListResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CommentListResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalCommentActionResponse(data []byte) (CommentActionResponse, error) {
	var r CommentActionResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CommentActionResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
