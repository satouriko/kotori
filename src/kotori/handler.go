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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(res_json)
	return
}

func Pong(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Pong!")
}

func GetComment(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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
	comments, err := ListComments(db, commentZoneID, fatherID, offsetID)
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

func StoreComment(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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
	comment.FatherID = fatherID
	comment.Type = "Comment"
	err = SaveComment(db, comment)
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
	}
	responseJson(w, res)
}

func DeleteComment(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	CommentID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred parsing comment id.", http.StatusInternalServerError)
		return
	}
	CommentID := uint(CommentID64)
	err = RemoveComment(db, CommentID)
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
	fmt.Fprint(w, "Login!")
}

func Logout(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Logout!")
}

func GetIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "GetIndex!")
}

func StoreIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "StoreIndex!")
}

func UpdateIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "UpdateIndex:"+ps.ByName("id"))
}

func DeleteIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "DeleteIndex:"+ps.ByName("id"))
}

func GetPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "GetPost!")
}

func StorePost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "StorePost!")
}

func UpdatePost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "UpdatePost:"+ps.ByName("id"))
}

func DeletePost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "DeletePost:"+ps.ByName("id"))
}
