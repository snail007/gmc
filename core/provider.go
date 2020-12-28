package gcore

import "io"

var Providers = NewProvider()

type SessionProvider func() Session
type SessionStorageProvider func(ctx Ctx) (SessionStorage, error)
type CacheProvider func(ctx Ctx) (Cache, error)
type TemplateProvider func(ctx Ctx) (Template, error)
type ViewProvider func(w io.Writer, tpl  Template) View
type DatabaseProvider func(ctx Ctx) (Database, error)
type DatabaseGroupProvider func(ctx Ctx) (DatabaseGroup, error)
type HTTPRouterProvider func(ctx Ctx) HTTPRouter
type ConfigProvider func() Config
type LoggerProvider func(ctx Ctx, prefix string) Logger

type ControllerProvider func() Controller
type CookiesProvider func(ctx Ctx) Cookies
type HTTPServerProvider func(ctx Ctx) HTTPServer
type APIServerProvider func(ctx Ctx) APIServer
type AppProvider func(ctx Ctx) App
type CtxProvider func() Ctx
type ErrorProvider func() Error
type I18nProvider func(ctx Ctx) (I18n, error)

type ProviderFactory struct {
	session        map[string]SessionProvider
	sessionStorage map[string]SessionStorageProvider
	cache          map[string]CacheProvider
	template       map[string]TemplateProvider
	view           map[string]ViewProvider
	database       map[string]DatabaseProvider
	databaseGroup  map[string]DatabaseGroupProvider
	httpRouter     map[string]HTTPRouterProvider
	config         map[string]ConfigProvider
	logger         map[string]LoggerProvider
	controller     map[string]ControllerProvider
	cookies        map[string]CookiesProvider
	httpServer     map[string]HTTPServerProvider
	apiServer      map[string]APIServerProvider
	app            map[string]AppProvider
	ctx            map[string]CtxProvider
	error          map[string]ErrorProvider
	i18n           map[string]I18nProvider
}

func (p *ProviderFactory) RegisterSession(key string, session SessionProvider) {
	p.session[key] = session
}

func (p *ProviderFactory) Session(key string) SessionProvider {
	return p.session[key]
}

func (p *ProviderFactory) SessionStorage(key string) SessionStorageProvider {
	return p.sessionStorage[key]
}

func (p *ProviderFactory) RegisterSessionStorage(key string, sessionStorage SessionStorageProvider) {
	p.sessionStorage[key] = sessionStorage
}

func (p *ProviderFactory) RegisterCache(key string, cache CacheProvider) {
	p.cache[key] = cache
}

func (p *ProviderFactory) Cache(key string) CacheProvider {
	return p.cache[key]
}

func (p *ProviderFactory) RegisterView(key string, view ViewProvider) {
	p.view[key] = view
}

func (p *ProviderFactory) View(key string) ViewProvider {
	return p.view[key]
}

func (p *ProviderFactory) RegisterTemplate(key string, template TemplateProvider) {
	p.template[key] = template
}

func (p *ProviderFactory) Template(key string) TemplateProvider {
	return p.template[key]
}

func (p *ProviderFactory) RegisterDatabase(key string, database DatabaseProvider) {
	p.database[key] = database
}

func (p *ProviderFactory) Database(key string) DatabaseProvider {
	return p.database[key]
}

func (p *ProviderFactory) RegisterDatabaseGroup(key string, databaseGroup DatabaseGroupProvider) {
	p.databaseGroup[key] = databaseGroup
}

func (p *ProviderFactory) DatabaseGroup(key string) DatabaseGroupProvider {
	return p.databaseGroup[key]
}

func (p *ProviderFactory) RegisterHTTPRouter(key string, router HTTPRouterProvider) {
	p.httpRouter[key] = router
}

func (p *ProviderFactory) HTTPRouter(key string) HTTPRouterProvider {
	return p.httpRouter[key]
}

func (p *ProviderFactory) RegisterConfig(key string, config ConfigProvider) {
	p.config[key] = config
}

func (p *ProviderFactory) Config(key string) ConfigProvider {
	return p.config[key]
}

func (p *ProviderFactory) RegisterLogger(key string, logger LoggerProvider) {
	p.logger[key] = logger
}

func (p *ProviderFactory) Logger(key string) LoggerProvider {
	return p.logger[key]
}

func (p *ProviderFactory) RegisterController(key string, controller ControllerProvider) {
	p.controller[key] = controller
}

func (p *ProviderFactory) Controller(key string) ControllerProvider {
	return p.controller[key]
}

func (p *ProviderFactory) RegisterCookies(key string, cookies CookiesProvider) {
	p.cookies[key] = cookies
}

func (p *ProviderFactory) Cookies(key string) CookiesProvider {
	return p.cookies[key]
}

func (p *ProviderFactory) RegisterHTTPServer(key string, httpServer HTTPServerProvider) {
	p.httpServer[key] = httpServer
}

func (p *ProviderFactory) HTTPServer(key string) HTTPServerProvider {
	return p.httpServer[key]
}

func (p *ProviderFactory) RegisterAPIServer(key string, apiServer APIServerProvider) {
	p.apiServer[key] = apiServer
}

func (p *ProviderFactory) APIServer(key string) APIServerProvider {
	return p.apiServer[key]
}

func (p *ProviderFactory) RegisterApp(key string, app AppProvider) {
	p.app[key] = app
}

func (p *ProviderFactory) App(key string) AppProvider {
	return p.app[key]
}

func (p *ProviderFactory) RegisterCtx(key string, ctx CtxProvider) {
	p.ctx[key] = ctx
}

func (p *ProviderFactory) Ctx(key string) CtxProvider {
	return p.ctx[key]
}

func (p *ProviderFactory) RegisterError(key string, error ErrorProvider) {
	p.error[key] = error
}

func (p *ProviderFactory) Error(key string) ErrorProvider {
	return p.error[key]
}

func (p *ProviderFactory) RegisterI18n(key string, i18n I18nProvider) {
	p.i18n[key] = i18n
}

func (p *ProviderFactory) I18n(key string) I18nProvider {
	return p.i18n[key]
}

func NewProvider() *ProviderFactory {
	return &ProviderFactory{
		session:        map[string]SessionProvider{},
		sessionStorage: map[string]SessionStorageProvider{},
		cache:          map[string]CacheProvider{},
		template:       map[string]TemplateProvider{},
		view:           map[string]ViewProvider{},
		database:       map[string]DatabaseProvider{},
		databaseGroup:  map[string]DatabaseGroupProvider{},
		httpRouter:     map[string]HTTPRouterProvider{},
		config:         map[string]ConfigProvider{},
		logger:         map[string]LoggerProvider{},
		controller:     map[string]ControllerProvider{},
		cookies:        map[string]CookiesProvider{},
		httpServer:     map[string]HTTPServerProvider{},
		apiServer:      map[string]APIServerProvider{},
		app:            map[string]AppProvider{},
		ctx:            map[string]CtxProvider{},
		error:          map[string]ErrorProvider{},
		i18n:           map[string]I18nProvider{},
	}
}
