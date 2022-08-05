package router

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"urlsortener/internal/infrastructure/api/openapi"
	"urlsortener/internal/usecases/app"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	http.Handler
	app *app.Urls
}


type URLRequest struct {
	OriginalUrl string  `json:"originalURL"`
}


type URLResponse struct {
	ShortUrl string	 `json:"shortURL"`
	StatisticsUrl string  `json:"statsURL"`
}

type Statistics struct{
	ShortUrl string  `json:"shortURL"`
	NumberRedirect int `json:"numberRedirect"`
}

func NewRouter(app *app.Urls) *Router{
	r := chi.NewRouter()
	ret := &Router{app: app}
	r.Use(middleware.Logger)

	r.Get("/", ret.GetMainPage)
	r.Get("/openapi", ret.GetOpenAPI)
	

	r.Mount("/", openapi.Handler(ret))

	swagger, err := openapi.GetSwagger()
	if err != nil {
		log.Fatalf("Error swagger: %v\n", err)
	}
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request){
		encoding := json.NewEncoder(w)
		_ = encoding.Encode(swagger)
	})
	
	ret.Handler = r
	return ret
}


func (ret *Router) CreateshortURL(w http.ResponseWriter, r *http.Request) {
	ru := &URLRequest{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(ru); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

	 _, err := url.ParseRequestURI(ru.OriginalUrl)
	 if err != nil{
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	url, err := ret.app.CreateURL(r.Context(), ru.OriginalUrl)
	if err != nil {
		http.Error(w, "Error when creating shor url", http.StatusBadRequest)
	}

	responsUrl := &URLResponse{
		ShortUrl: "/" + url.ShortUrl,
		StatisticsUrl: "/statistic" + url.ShortUrl,
	}

	w.Header().Add("Content-type", "applocation/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(responsUrl)
}

func (ret *Router) RedirectURL(w http.ResponseWriter, r *http.Request, shortUrl string) {
	url, err := ret.app.RedirectUrl(r.Context(), shortUrl)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusSeeOther)
}

func (ret *Router) GetStats(w http.ResponseWriter, r *http.Request, shortUrl string) {
	s, err := ret.app.Stats(r.Context(), shortUrl)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	resp := &Statistics{
		ShortUrl: s.ShortURL,
		NumberRedirect: s.NumberRedirect,
	}

	w.Header().Add("Content-type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func(ret *Router) GetMainPage(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("./static/index.html")
	if err != nil{ 
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err != temp.Execute(w, nil) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

}

func(ret *Router) GetOpenAPI(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles()
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	if err != t.Execute(w, nil) {
		http.Error(w, "Server error", 500)
		return
	}

}
