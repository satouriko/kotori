package kotori

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/yanzay/log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func respondJson(w http.ResponseWriter, data map[string]interface{}, httpStatusCode int) {
	resJson, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred encoding response.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	w.Write(resJson)
	return
}

func checkAdmin(w http.ResponseWriter, req *http.Request) (result bool) {
	sess, _ := globalSessions.SessionStart(w, req)
	defer sess.SessionRelease(w)
	if priv := sess.Get("privilege"); priv == nil || priv.(string) != "admin" {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    "Authorize failed.",
		}
		respondJson(w, res, http.StatusUnauthorized)
		result = false
		return
	}
	result = true
	return
}

func Pong(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"msg":    "Pong!",
	}
	respondJson(w, res, http.StatusOK)
}

func Status(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"uptime": time.Since(startTime).String(),
		"d":      int(time.Since(startTime).Hours() / 24),
	}
	respondJson(w, res, http.StatusUnauthorized)
	return
}

func ListComment(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	if len(req.Form["comment_zone_id"]) != 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid comment zone.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	commentZoneID64, err := strconv.ParseUint(req.Form["comment_zone_id"][0], 10, 32)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing comment zone id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	commentZoneID := uint(commentZoneID64)
	count, err := CountComments(db, commentZoneID)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred querying comments.",
		}
		respondJson(w, res, http.StatusInternalServerError)
		return
	}
	if len(req.Form["count"]) == 1 {
		if req.Form["count"][0] != "" {
			res := map[string]interface{}{
				"code":   http.StatusOK,
				"result": true,
				"cnt":    count,
			}
			respondJson(w, res, http.StatusOK)
		}
		return
	}
	var fatherID uint
	if len(req.Form["father_id"]) > 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid father id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	} else if len(req.Form["father_id"]) == 1 {
		fatherID64, err := strconv.ParseUint(req.Form["father_id"][0], 10, 32)
		if err != nil {
			log.Error(err)
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Error occurred parsing father id.",
			}
			respondJson(w, res, http.StatusBadRequest)
			return
		}
		fatherID = uint(fatherID64)
	} else {
		fatherID = 0
	}
	var offsetID uint
	if len(req.Form["offset_id"]) > 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid offset id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	} else if len(req.Form["offset_id"]) == 1 {
		offsetID64, err := strconv.ParseUint(req.Form["offset_id"][0], 10, 32)
		if err != nil {
			log.Error(err)
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Error occurred parsing offset id.",
			}
			respondJson(w, res, http.StatusBadRequest)
			return
		}
		offsetID = uint(offsetID64)
	} else {
		offsetID = 0
	}
	comments, err := FindComments(db, commentZoneID, fatherID, offsetID)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred querying comments.",
		}
		respondJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   comments,
		"cnt":    count,
	}
	respondJson(w, res, http.StatusOK)
}

func CreateComment(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	var comment Comment
	if len(req.Form["comment_zone_id"]) != 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid comment zone.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	commentZoneID64, err := strconv.ParseUint(req.Form["comment_zone_id"][0], 10, 32)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing comment zone id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	commentZoneID := uint(commentZoneID64)
	comment.CommentZoneID = commentZoneID
	if len(req.Form["content"]) != 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid comment content.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	comment.Content = req.Form["content"][0]
	if len(req.Form["name"]) != 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid user name.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	comment.User.Name = req.Form["name"][0]
	if len(req.Form["email"]) != 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid user email.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	comment.User.Email = req.Form["email"][0]
	if len(req.Form["website"]) > 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid user website.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	} else if len(req.Form["website"]) == 1 {
		comment.User.Website = req.Form["website"][0]
	}
	var fatherID uint
	if len(req.Form["father_id"]) > 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid father id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	} else if len(req.Form["father_id"]) == 1 {
		fatherID64, err := strconv.ParseUint(req.Form["father_id"][0], 10, 32)
		if err != nil {
			log.Error(err)
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Error occurred parsing father id.",
			}
			respondJson(w, res, http.StatusBadRequest)
			return
		}
		fatherID = uint(fatherID64)
	} else {
		fatherID = 0
	}
	var replyUserID uint
	if len(req.Form["reply_user_id"]) > 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid reply user id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	} else if len(req.Form["reply_user_id"]) == 1 {
		replyUserID64, err := strconv.ParseUint(req.Form["reply_user_id"][0], 10, 32)
		if err != nil {
			log.Error(err)
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Error occurred parsing reply user id.",
			}
			respondJson(w, res, http.StatusBadRequest)
			return
		}
		replyUserID = uint(replyUserID64)
	} else {
		replyUserID = 0
	}
	comment.FatherID = fatherID
	comment.ReplyUserID = replyUserID
	comment.Type = "Comment"
	comment, err = StoreComment(db, comment)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred storing comment to database.",
		}
		respondJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   comment,
	}
	respondJson(w, res, http.StatusOK)
}

func DeleteComment(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !checkAdmin(w, req) {
		return
	}

	commentID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing comment id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	commentID := uint(commentID64)
	err = RemoveComment(db, commentID)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred removing comment from database.",
		}
		respondJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
	}
	respondJson(w, res, http.StatusOK)
}

func Login(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	sess, _ := globalSessions.SessionStart(w, req)
	defer sess.SessionRelease(w)
	if username := sess.Get("username"); username != nil {
		res := map[string]interface{}{
			"code":   http.StatusOK,
			"result": false,
			"msg":    "Already logged in as: " + username.(string),
		}
		respondJson(w, res, http.StatusOK)
	} else {
		req.ParseForm()
		if len(req.Form["username"]) != 1 {
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Invalid username.",
			}
			respondJson(w, res, http.StatusBadRequest)
			return
		}
		username = req.Form["username"][0]
		if len(req.Form["password"]) != 1 {
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Invalid password.",
			}
			respondJson(w, res, http.StatusBadRequest)
			return
		}
		password := req.Form["password"][0]
		for _, admin := range GlobCfg.ADMIN {
			if admin.Username == username && admin.Password == password {
				sess.Set("username", username)
				sess.Set("privilege", "admin")
				res := map[string]interface{}{
					"code":   http.StatusOK,
					"result": true,
					"msg":    "Successfully logged in as: " + username.(string),
				}
				respondJson(w, res, http.StatusOK)
				return
			}
		}
		res := map[string]interface{}{
			"code":   http.StatusOK,
			"result": false,
			"msg":    "No such user or password mismatch.",
		}
		respondJson(w, res, http.StatusOK)
		return
	}
	log.Info(GlobCfg.ADMIN)
}

func Logout(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	sess, _ := globalSessions.SessionStart(w, req)
	defer sess.SessionRelease(w)
	sess.Delete("username")
	sess.Delete("privilege")
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"msg":    "Successfully logged out.",
	}
	respondJson(w, res, http.StatusOK)
}

func EditUserSetHonor(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !checkAdmin(w, req) {
		return
	}

	req.ParseForm()
	userID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing user id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	userID := uint(userID64)
	if len(req.Form["honor"]) != 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid honor.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	honor := req.Form["honor"][0]
	user, err := UpdateUserSetHonor(db, userID, honor)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred storing user to database: " + err.Error(),
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   user,
	}
	respondJson(w, res, http.StatusOK)
}

func ListIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	if len(req.Form["class"]) != 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid index class.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	class := req.Form["class"][0]
	var offsetID uint
	if len(req.Form["offset_id"]) > 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid offset id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	} else if len(req.Form["offset_id"]) == 1 {
		offsetID64, err := strconv.ParseUint(req.Form["offset_id"][0], 10, 32)
		if err != nil {
			log.Error(err)
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Error occurred parsing offset id.",
			}
			respondJson(w, res, http.StatusBadRequest)
			return
		}
		offsetID = uint(offsetID64)
	} else {
		offsetID = 0
	}
	var order = "asc"
	if len(req.Form["order"]) == 1 && req.Form["order"][0] != "asc" {
		order = "desc"
	}
	indexes, err := FindIndexes(db, class, order, offsetID)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred querying indexes.",
		}
		respondJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   indexes,
	}
	respondJson(w, res, http.StatusOK)
}

func GetIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var index Index
	var err error
	if req.Header.Get("X-Query-By") == "Title" {
		indexTitle := ps.ByName("id")
		index, err = FindIndexByTitle(db, indexTitle)
	} else {
		indexID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
		if err != nil {
			log.Error(err)
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Error occurred parsing index id.",
			}
			respondJson(w, res, http.StatusBadRequest)
			return
		}
		indexID := uint(indexID64)
		index, err = FindIndex(db, indexID)
	}
	if err != nil {
		log.Error(err)
		if strings.Contains(err.Error(), "record not found") {
			res := map[string]interface{}{
				"code":   http.StatusNotFound,
				"result": false,
				"msg":    "Index not found.",
			}
			respondJson(w, res, http.StatusNotFound)
			return
		}
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred querying index from database.",
		}
		respondJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   index,
	}
	respondJson(w, res, http.StatusOK)
}

func CreateIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !checkAdmin(w, req) {
		return
	}

	req.ParseForm()
	var index Index
	if len(req.Form["class"]) != 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid index class.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	index.Class = req.Form["class"][0]
	if len(req.Form["attr"]) == 1 {
		index.Attr = req.Form["attr"][0]
	}
	if len(req.Form["title"]) == 1 {
		index.Title = req.Form["title"][0]
	}
	index, err := StoreIndex(db, index)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred storing index to database.",
		}
		respondJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   index,
	}
	respondJson(w, res, http.StatusOK)
}

func EditIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !checkAdmin(w, req) {
		return
	}

	req.ParseForm()
	var index Index
	indexID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing index id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	indexID := uint(indexID64)
	index.ID = indexID
	if len(req.Form["class"]) != 0 {
		res := map[string]interface{}{
			"code":   http.StatusForbidden,
			"result": false,
			"msg":    "Class could not be changed.",
		}
		respondJson(w, res, http.StatusForbidden)
		return
	}
	if len(req.Form["attr"]) == 1 {
		index.Attr = req.Form["attr"][0]
	}
	if len(req.Form["title"]) == 1 {
		index.Title = req.Form["title"][0]
	}
	index, err = UpdateIndex(db, index)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred storing index to database.",
		}
		respondJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   index,
	}
	respondJson(w, res, http.StatusOK)
}

func DeleteIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !checkAdmin(w, req) {
		return
	}

	indexID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing index id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	indexID := uint(indexID64)
	err = RemoveIndex(db, indexID)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred removing index from database.",
		}
		respondJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
	}
	respondJson(w, res, http.StatusOK)
}

func ListPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	var offsetID uint
	if len(req.Form["offset_id"]) > 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid offset id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	} else if len(req.Form["offset_id"]) == 1 {
		offsetID64, err := strconv.ParseUint(req.Form["offset_id"][0], 10, 32)
		if err != nil {
			log.Error(err)
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Error occurred parsing offset id.",
			}
			respondJson(w, res, http.StatusBadRequest)
			return
		}
		offsetID = uint(offsetID64)
	} else {
		offsetID = 0
	}
	posts, err := FindPosts(db, offsetID)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred querying posts.",
		}
		respondJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   posts,
	}
	respondJson(w, res, http.StatusOK)
}

func GetPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	postID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing post id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	postID := uint(postID64)
	post, err := FindPost(db, postID)
	if err != nil {
		log.Error(err)
		if strings.Contains(err.Error(), "record not found") {
			res := map[string]interface{}{
				"code":   http.StatusNotFound,
				"result": false,
				"msg":    "Post not found.",
			}
			respondJson(w, res, http.StatusNotFound)
			return
		}
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred querying post from database.",
		}
		respondJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   post,
	}
	respondJson(w, res, http.StatusOK)
}

func CreatePost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !checkAdmin(w, req) {
		return
	}

	req.ParseForm()
	var post Post
	if len(req.Form["content"]) == 1 {
		post.Content = req.Form["content"][0]
	}
	if len(req.Form["title"]) == 1 {
		post.Title = req.Form["title"][0]
	}
	post, err := StorePost(db, post)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred storing post to database.",
		}
		respondJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   post,
	}
	respondJson(w, res, http.StatusOK)
}

func EditPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !checkAdmin(w, req) {
		return
	}

	req.ParseForm()
	var post Post
	postID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing index id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	postID := uint(postID64)
	post.ID = postID
	if len(req.Form["content"]) == 1 {
		post.Content = req.Form["content"][0]
	}
	if len(req.Form["title"]) == 1 {
		post.Title = req.Form["title"][0]
	}
	post, err = UpdatePost(db, post)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred storing post to database.",
		}
		respondJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   post,
	}
	respondJson(w, res, http.StatusOK)
}

func DeletePost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !checkAdmin(w, req) {
		return
	}

	postID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing index id.",
		}
		respondJson(w, res, http.StatusBadRequest)
		return
	}
	postID := uint(postID64)
	err = RemovePost(db, postID)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred removing post from database.",
		}
		respondJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
	}
	respondJson(w, res, http.StatusOK)
}
