package main

import (
	"net/http"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"strconv"
	"github.com/yanzay/log"
	"encoding/json"
)

func responseJson(w http.ResponseWriter, data map[string]interface{}) {
	res_json, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred encoding response.", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res_json)
	return
}

func checkAdmin(w http.ResponseWriter, req *http.Request) (result bool)  {
	sess, _ := globalSessions.SessionStart(w, req)
	defer sess.SessionRelease(w)
	if priv := sess.Get("privilege"); priv == nil || priv.(string) != "admin" {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Authorize failed.",
		}
		responseJson(w, res)
		result = false;
		return
	}
	result = true
	return
}

func Pong(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Pong!")
}

func ListComment(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	if (len(req.Form["comment_zone_id"]) != 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid comment zone.",
		}
		responseJson(w, res)
		return
	}
	commentZoneID64, err := strconv.ParseUint(req.Form["comment_zone_id"][0], 10, 32)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred parsing comment zone id.", http.StatusInternalServerError)
		return
	}
	commentZoneID := uint(commentZoneID64)
	var fatherID uint
	if (len(req.Form["father_id"]) > 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid father id.",
		}
		responseJson(w, res)
		return
	} else if (len(req.Form["father_id"]) == 1) {
		fatherID64, err := strconv.ParseUint(req.Form["father_id"][0], 10, 32)
		if err != nil {
			log.Error(err)
			http.Error(w, "Error occurred parsing father id.", http.StatusInternalServerError)
			return
		}
		fatherID = uint(fatherID64)
	} else {
		fatherID = 0
	}
	var offsetID uint
	if (len(req.Form["offset_id"]) > 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid offset id.",
		}
		responseJson(w, res)
		return
	} else if (len(req.Form["offset_id"]) == 1) {
		offsetID64, err := strconv.ParseUint(req.Form["offset_id"][0], 10, 32)
		if err != nil {
			log.Error(err)
			http.Error(w, "Error occurred parsing offset id.", http.StatusInternalServerError)
			return
		}
		offsetID = uint(offsetID64)
	} else {
		offsetID = 0
	}
	comments, err := FindComments(db, commentZoneID, fatherID, offsetID)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred querying comments.", http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"result": true,
		"data":   comments,
	}
	responseJson(w, res)
}

func CreateComment(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	var comment Comment
	if (len(req.Form["comment_zone_id"]) != 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid comment zone.",
		}
		responseJson(w, res)
		return
	}
	commentZoneID64, err := strconv.ParseUint(req.Form["comment_zone_id"][0], 10, 32)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred parsing comment zone id.", http.StatusInternalServerError)
		return
	}
	commentZoneID := uint(commentZoneID64)
	comment.CommentZoneID = commentZoneID
	if (len(req.Form["content"]) != 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid comment content.",
		}
		responseJson(w, res)
		return
	}
	comment.Content = req.Form["content"][0]
	if (len(req.Form["name"]) != 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid user name.",
		}
		responseJson(w, res)
		return
	}
	comment.User.Name = req.Form["name"][0]
	if (len(req.Form["email"]) != 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid user email.",
		}
		responseJson(w, res)
		return
	}
	comment.User.Email = req.Form["email"][0]
	if (len(req.Form["website"]) > 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid user website.",
		}
		responseJson(w, res)
		return
	} else if (len(req.Form["website"]) == 1) {
		comment.User.Website = req.Form["website"][0]
	}
	var fatherID uint
	if (len(req.Form["father_id"]) > 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid father id.",
		}
		responseJson(w, res)
		return
	} else if (len(req.Form["father_id"]) == 1) {
		fatherID64, err := strconv.ParseUint(req.Form["father_id"][0], 10, 32)
		if err != nil {
			log.Error(err)
			http.Error(w, "Error occurred parsing father id.", http.StatusInternalServerError)
			return
		}
		fatherID = uint(fatherID64)
	} else {
		fatherID = 0
	}
	var replyUserID uint
	if (len(req.Form["reply_user_id"]) > 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid reply user id.",
		}
		responseJson(w, res)
		return
	} else if (len(req.Form["reply_user_id"]) == 1) {
		replyUserID64, err := strconv.ParseUint(req.Form["reply_user_id"][0], 10, 32)
		if err != nil {
			log.Error(err)
			http.Error(w, "Error occurred parsing reply user id.", http.StatusInternalServerError)
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
			"result": false,
			"msg":    "Error occurred storing comment to database: " + err.Error(),
		}
		responseJson(w, res)
		return
	}
	res := map[string]interface{}{
		"result": true,
		"data": comment,
	}
	responseJson(w, res)
}

func DeleteComment(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !checkAdmin(w, req) {
		return
	}

	commentID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred parsing comment id.", http.StatusInternalServerError)
		return
	}
	commentID := uint(commentID64)
	err = RemoveComment(db, commentID)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"result": false,
			"msg":    "Error occurred removing comment from database: " + err.Error(),
		}
		responseJson(w, res)
		return
	}
	res := map[string]interface{}{
		"result": true,
	}
	responseJson(w, res)
}

func Login(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	sess, _ := globalSessions.SessionStart(w, req)
	defer sess.SessionRelease(w)
	if username := sess.Get("username"); username != nil {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Already logged in as: " + username.(string),
		}
		responseJson(w, res)
	} else {
		req.ParseForm()
		if (len(req.Form["username"]) != 1) {
			res := map[string]interface{}{
				"result": false,
				"msg":    "Invalid username.",
			}
			responseJson(w, res)
			return
		}
		username = req.Form["username"][0]
		if (len(req.Form["password"]) != 1) {
			res := map[string]interface{}{
				"result": false,
				"msg":    "Invalid password.",
			}
			responseJson(w, res)
			return
		}
		password := req.Form["password"][0]
		for _, admin := range GlobCfg.ADMIN {
			if admin.Username == username && admin.Password == password {
				sess.Set("username", username)
				sess.Set("privilege", "admin")
				res := map[string]interface{}{
					"result": true,
					"msg":    "Successfully logged in as: " + username.(string),
				}
				responseJson(w, res)
				return
			}
		}
		res := map[string]interface{}{
			"result": false,
			"msg":    "No such user or password mismatch.",
		}
		responseJson(w, res)
		return
	}
	fmt.Print(GlobCfg.ADMIN)
}

func Logout(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	sess, _ := globalSessions.SessionStart(w, req)
	defer sess.SessionRelease(w)
	sess.Delete("username")
	sess.Delete("privilege")
	res := map[string]interface{}{
		"result": true,
		"msg":    "Successfully logged out.",
	}
	responseJson(w, res)
}

func EditUserSetHonor(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !checkAdmin(w, req) {
		return
	}

	req.ParseForm()
	userID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred parsing user id.", http.StatusInternalServerError)
		return
	}
	userID := uint(userID64)
	if (len(req.Form["honor"]) != 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid honor.",
		}
		responseJson(w, res)
		return
	}
	honor := req.Form["honor"][0]
	user, err := UpdateUserSetHonor(db, userID, honor)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"result": false,
			"msg":    "Error occurred storing user to database: " + err.Error(),
		}
		responseJson(w, res)
		return
	}
	res := map[string]interface{}{
		"result": true,
		"data": user,
	}
	responseJson(w, res)
}

func ListIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	if (len(req.Form["class"]) != 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid index class.",
		}
		responseJson(w, res)
		return
	}
	class := req.Form["class"][0]
	var offsetID uint
	if (len(req.Form["offset_id"]) > 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid offset id.",
		}
		responseJson(w, res)
		return
	} else if (len(req.Form["offset_id"]) == 1) {
		offsetID64, err := strconv.ParseUint(req.Form["offset_id"][0], 10, 32)
		if err != nil {
			log.Error(err)
			http.Error(w, "Error occurred parsing offset id.", http.StatusInternalServerError)
			return
		}
		offsetID = uint(offsetID64)
	} else {
		offsetID = 0
	}
	var order = "asc"
	if (len(req.Form["order"]) == 1 && req.Form["order"][0] != "asc") {
		order = "desc"
	}
	indexes, err := FindIndexes(db, class, order, offsetID)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred querying indexes.", http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"result": true,
		"data":   indexes,
	}
	responseJson(w, res)
}

func CreateIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !checkAdmin(w, req) {
		return
	}

	req.ParseForm()
	var index Index
	if (len(req.Form["class"]) != 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid index class.",
		}
		responseJson(w, res)
		return
	}
	index.Class = req.Form["class"][0]
	if (len(req.Form["attr"]) == 1) {
		index.Attr = req.Form["attr"][0]
	}
	if (len(req.Form["title"]) == 1) {
		index.Title = req.Form["title"][0]
	}
	index, err := StoreIndex(db, index)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"result": false,
			"msg":    "Error occurred storing index to database: " + err.Error(),
		}
		responseJson(w, res)
		return
	}
	res := map[string]interface{}{
		"result": true,
		"data": index,
	}
	responseJson(w, res)
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
		http.Error(w, "Error occurred parsing index id.", http.StatusInternalServerError)
		return
	}
	indexID := uint(indexID64)
	index.ID = indexID
	if (len(req.Form["class"]) != 0) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Class could not be changed.",
		}
		responseJson(w, res)
		return
	}
	if (len(req.Form["attr"]) == 1) {
		index.Attr = req.Form["attr"][0]
	}
	if (len(req.Form["title"]) == 1) {
		index.Title = req.Form["title"][0]
	}
	index, err = UpdateIndex(db, index)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"result": false,
			"msg":    "Error occurred storing index to database: " + err.Error(),
		}
		responseJson(w, res)
		return
	}
	res := map[string]interface{}{
		"result": true,
		"data": index,
	}
	responseJson(w, res)
}

func DeleteIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !checkAdmin(w, req) {
		return
	}

	indexID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred parsing index id.", http.StatusInternalServerError)
		return
	}
	indexID := uint(indexID64)
	err = RemoveIndex(db, indexID)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"result": false,
			"msg":    "Error occurred removing index from database: " + err.Error(),
		}
		responseJson(w, res)
		return
	}
	res := map[string]interface{}{
		"result": true,
	}
	responseJson(w, res)
}

func ListPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	var offsetID uint
	if (len(req.Form["offset_id"]) > 1) {
		res := map[string]interface{}{
			"result": false,
			"msg":    "Invalid offset id.",
		}
		responseJson(w, res)
		return
	} else if (len(req.Form["offset_id"]) == 1) {
		offsetID64, err := strconv.ParseUint(req.Form["offset_id"][0], 10, 32)
		if err != nil {
			log.Error(err)
			http.Error(w, "Error occurred parsing offset id.", http.StatusInternalServerError)
			return
		}
		offsetID = uint(offsetID64)
	} else {
		offsetID = 0
	}
	posts, err := FindPosts(db, offsetID)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred querying posts.", http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"result": true,
		"data":   posts,
	}
	responseJson(w, res)
}

func GetPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	postID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred parsing post id.", http.StatusInternalServerError)
		return
	}
	postID := uint(postID64)
	post, err := FindPost(db, postID)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"result": false,
			"msg":    "Error occurred querying post from database: " + err.Error(),
		}
		responseJson(w, res)
		return
	}
	res := map[string]interface{}{
		"result": true,
		"data": post,
	}
	responseJson(w, res)
}

func CreatePost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !checkAdmin(w, req) {
		return
	}

	req.ParseForm()
	var post Post
	if (len(req.Form["content"]) == 1) {
		post.Content = req.Form["content"][0]
	}
	if (len(req.Form["title"]) == 1) {
		post.Title = req.Form["title"][0]
	}
	post, err := StorePost(db, post)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"result": false,
			"msg":    "Error occurred storing post to database: " + err.Error(),
		}
		responseJson(w, res)
		return
	}
	res := map[string]interface{}{
		"result": true,
		"data": post,
	}
	responseJson(w, res)
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
		http.Error(w, "Error occurred parsing index id.", http.StatusInternalServerError)
		return
	}
	postID := uint(postID64)
	post.ID = postID
	if (len(req.Form["content"]) == 1) {
		post.Content = req.Form["content"][0]
	}
	if (len(req.Form["title"]) == 1) {
		post.Title = req.Form["title"][0]
	}
	post, err = UpdatePost(db, post)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"result": false,
			"msg":    "Error occurred storing post to database: " + err.Error(),
		}
		responseJson(w, res)
		return
	}
	res := map[string]interface{}{
		"result": true,
		"data": post,
	}
	responseJson(w, res)
}

func DeletePost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !checkAdmin(w, req) {
		return
	}

	postID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred parsing index id.", http.StatusInternalServerError)
		return
	}
	postID := uint(postID64)
	err = RemovePost(db, postID)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"result": false,
			"msg":    "Error occurred removing post from database: " + err.Error(),
		}
		responseJson(w, res)
		return
	}
	res := map[string]interface{}{
		"result": true,
	}
	responseJson(w, res)
}
