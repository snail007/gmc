// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gcore

import "io"

const (
	DefaultProviderKey = "default"
)

var defaultAutoProvider = NewAutoProvider()

func RegisterSession(key string, session SessionProvider) {
	defaultAutoProvider.RegisterSession(key, session)
}

func ProviderSession(key ...string) SessionProvider {
	return defaultAutoProvider.Session(key...)
}

func RegisterSessionStorage(key string, sessionStorage SessionStorageProvider) {
	defaultAutoProvider.RegisterSessionStorage(key, sessionStorage)
}

func ProviderSessionStorage(key ...string) SessionStorageProvider {
	return defaultAutoProvider.SessionStorage(key...)
}

func RegisterCache(key string, cache CacheProvider) {
	defaultAutoProvider.RegisterCache(key, cache)
}

func ProviderCache(key ...string) CacheProvider {
	return defaultAutoProvider.Cache(key...)
}

func RegisterView(key string, view ViewProvider) {
	defaultAutoProvider.RegisterView(key, view)
}

func ProviderView(key ...string) ViewProvider {
	return defaultAutoProvider.View(key...)
}

func RegisterTemplate(key string, template TemplateProvider) {
	defaultAutoProvider.RegisterTemplate(key, template)
}

func ProviderTemplate(key ...string) TemplateProvider {
	return defaultAutoProvider.Template(key...)
}

func RegisterDatabase(key string, database DatabaseProvider) {
	defaultAutoProvider.RegisterDatabase(key, database)
}

func ProviderDatabase(key ...string) DatabaseProvider {
	return defaultAutoProvider.Database(key...)
}

func RegisterDatabaseGroup(key string, databaseGroup DatabaseGroupProvider) {
	defaultAutoProvider.RegisterDatabaseGroup(key, databaseGroup)
}

func ProviderDatabaseGroup(key ...string) DatabaseGroupProvider {
	return defaultAutoProvider.DatabaseGroup(key...)
}

func RegisterHTTPRouter(key string, router HTTPRouterProvider) {
	defaultAutoProvider.RegisterHTTPRouter(key, router)
}

func ProviderHTTPRouter(key ...string) HTTPRouterProvider {
	return defaultAutoProvider.HTTPRouter(key...)
}

func RegisterConfig(key string, config ConfigProvider) {
	defaultAutoProvider.RegisterConfig(key, config)
}

func ProviderConfig(key ...string) ConfigProvider {
	return defaultAutoProvider.Config(key...)
}

func RegisterLogger(key string, logger LoggerProvider) {
	defaultAutoProvider.RegisterLogger(key, logger)
}

func ProviderLogger(key ...string) LoggerProvider {
	return defaultAutoProvider.Logger(key...)
}

func RegisterController(key string, controller ControllerProvider) {
	defaultAutoProvider.RegisterController(key, controller)
}

func ProviderController(key ...string) ControllerProvider {
	return defaultAutoProvider.Controller(key...)
}

func RegisterCookies(key string, cookies CookiesProvider) {
	defaultAutoProvider.RegisterCookies(key, cookies)
}

func ProviderCookies(key ...string) CookiesProvider {
	return defaultAutoProvider.Cookies(key...)
}

func RegisterHTTPServer(key string, httpServer HTTPServerProvider) {
	defaultAutoProvider.RegisterHTTPServer(key, httpServer)
}

func ProviderHTTPServer(key ...string) HTTPServerProvider {
	return defaultAutoProvider.HTTPServer(key...)
}

func RegisterAPIServer(key string, apiServer APIServerProvider) {
	defaultAutoProvider.RegisterAPIServer(key, apiServer)
}

func ProviderAPIServer(key ...string) APIServerProvider {
	return defaultAutoProvider.APIServer(key...)
}

func RegisterApp(key string, app AppProvider) {
	defaultAutoProvider.RegisterApp(key, app)
}

func ProviderApp(key ...string) AppProvider {
	return defaultAutoProvider.App(key...)
}

func RegisterCtx(key string, ctx CtxProvider) {
	defaultAutoProvider.RegisterCtx(key, ctx)
}

func ProviderCtx(key ...string) CtxProvider {
	return defaultAutoProvider.Ctx(key...)
}

func RegisterError(key string, error ErrorProvider) {
	defaultAutoProvider.RegisterError(key, error)
}

func ProviderError(key ...string) ErrorProvider {
	return defaultAutoProvider.Error(key...)
}

func RegisterI18n(key string, i18n I18nProvider) {
	defaultAutoProvider.RegisterI18n(key, i18n)
}

func ProviderI18n(key ...string) I18nProvider {
	return defaultAutoProvider.I18n(key...)
}

type AutoProvider struct {
	sessionKeys        []string
	sessionStorageKeys []string
	cacheKeys          []string
	templateKeys       []string
	viewKeys           []string
	databaseKeys       []string
	databaseGroupKeys  []string
	httpRouterKeys     []string
	configKeys         []string
	loggerKeys         []string
	controllerKeys     []string
	cookiesKeys        []string
	httpServerKeys     []string
	apiServerKeys      []string
	appKeys            []string
	ctxKeys            []string
	errorKeys          []string
	i18nKeys           []string
	factory            *ProviderFactory
}

func (p *AutoProvider) find(providerKeys []string, callback func(key string) interface{}) interface{} {
	for i := len(providerKeys) - 1; i >= 0; i-- {
		if v := callback(providerKeys[i]); v != nil {
			return v
		}
	}
	return callback(DefaultProviderKey)
}

func (p *AutoProvider) addKey(providerKeys *[]string, key string) {
	found := false
	for _, v := range *providerKeys {
		if v == key {
			found = true
			break
		}
	}
	if !found {
		*providerKeys = append(*providerKeys, key)
	}
}

func (p *AutoProvider) RegisterSession(key string, session SessionProvider) {
	p.addKey(&p.sessionKeys, key)
	p.factory.RegisterSession(key, session)
}

func (p *AutoProvider) Session(key ...string) SessionProvider {
	if len(key) == 1 {
		return p.factory.Session(key[0])
	}
	if v := p.find(p.sessionKeys, func(key string) interface{} {
		return p.factory.Session(key)
	}); v != nil {
		return v.(SessionProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterSessionStorage(key string, sessionStorage SessionStorageProvider) {
	p.addKey(&p.sessionStorageKeys, key)
	p.factory.RegisterSessionStorage(key, sessionStorage)
}

func (p *AutoProvider) SessionStorage(key ...string) SessionStorageProvider {
	if len(key) == 1 {
		return p.factory.SessionStorage(key[0])
	}
	if v := p.find(p.sessionStorageKeys, func(key string) interface{} {
		return p.factory.SessionStorage(key)
	}); v != nil {
		return v.(SessionStorageProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterCache(key string, cache CacheProvider) {
	p.addKey(&p.cacheKeys, key)
	p.factory.RegisterCache(key, cache)
}

func (p *AutoProvider) Cache(key ...string) CacheProvider {
	if len(key) == 1 {
		return p.factory.Cache(key[0])
	}
	if v := p.find(p.cacheKeys, func(key string) interface{} {
		return p.factory.Cache(key)
	}); v != nil {
		return v.(CacheProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterView(key string, view ViewProvider) {
	p.addKey(&p.viewKeys, key)
	p.factory.RegisterView(key, view)
}

func (p *AutoProvider) View(key ...string) ViewProvider {
	if len(key) == 1 {
		return p.factory.View(key[0])
	}
	if v := p.find(p.viewKeys, func(key string) interface{} {
		return p.factory.View(key)
	}); v != nil {
		return v.(ViewProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterTemplate(key string, template TemplateProvider) {
	p.addKey(&p.templateKeys, key)
	p.factory.RegisterTemplate(key, template)
}

func (p *AutoProvider) Template(key ...string) TemplateProvider {
	if len(key) == 1 {
		return p.factory.Template(key[0])
	}
	if v := p.find(p.templateKeys, func(key string) interface{} {
		return p.factory.Template(key)
	}); v != nil {
		return v.(TemplateProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterDatabase(key string, database DatabaseProvider) {
	p.addKey(&p.databaseKeys, key)
	p.factory.RegisterDatabase(key, database)
}

func (p *AutoProvider) Database(key ...string) DatabaseProvider {
	if len(key) == 1 {
		return p.factory.Database(key[0])
	}
	if v := p.find(p.databaseKeys, func(key string) interface{} {
		return p.factory.Database(key)
	}); v != nil {
		return v.(DatabaseProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterDatabaseGroup(key string, databaseGroup DatabaseGroupProvider) {
	p.addKey(&p.databaseGroupKeys, key)
	p.factory.RegisterDatabaseGroup(key, databaseGroup)
}

func (p *AutoProvider) DatabaseGroup(key ...string) DatabaseGroupProvider {
	if len(key) == 1 {
		return p.factory.DatabaseGroup(key[0])
	}
	if v := p.find(p.databaseGroupKeys, func(key string) interface{} {
		return p.factory.DatabaseGroup(key)
	}); v != nil {
		return v.(DatabaseGroupProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterHTTPRouter(key string, router HTTPRouterProvider) {
	p.addKey(&p.httpRouterKeys, key)
	p.factory.RegisterHTTPRouter(key, router)
}

func (p *AutoProvider) HTTPRouter(key ...string) HTTPRouterProvider {
	if len(key) == 1 {
		return p.factory.HTTPRouter(key[0])
	}
	if v := p.find(p.httpRouterKeys, func(key string) interface{} {
		return p.factory.HTTPRouter(key)
	}); v != nil {
		return v.(HTTPRouterProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterConfig(key string, config ConfigProvider) {
	p.addKey(&p.configKeys, key)
	p.factory.RegisterConfig(key, config)
}

func (p *AutoProvider) Config(key ...string) ConfigProvider {
	if len(key) == 1 {
		return p.factory.Config(key[0])
	}
	if v := p.find(p.configKeys, func(key string) interface{} {
		return p.factory.Config(key)
	}); v != nil {
		return v.(ConfigProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterLogger(key string, logger LoggerProvider) {
	p.addKey(&p.loggerKeys, key)
	p.factory.RegisterLogger(key, logger)
}

func (p *AutoProvider) Logger(key ...string) LoggerProvider {
	if len(key) == 1 {
		return p.factory.Logger(key[0])
	}
	if v := p.find(p.loggerKeys, func(key string) interface{} {
		return p.factory.Logger(key)
	}); v != nil {
		return v.(LoggerProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterController(key string, controller ControllerProvider) {
	p.addKey(&p.controllerKeys, key)
	p.factory.RegisterController(key, controller)
}

func (p *AutoProvider) Controller(key ...string) ControllerProvider {
	if len(key) == 1 {
		return p.factory.Controller(key[0])
	}
	if v := p.find(p.controllerKeys, func(key string) interface{} {
		return p.factory.Controller(key)
	}); v != nil {
		return v.(ControllerProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterCookies(key string, cookies CookiesProvider) {
	p.addKey(&p.cookiesKeys, key)
	p.factory.RegisterCookies(key, cookies)
}

func (p *AutoProvider) Cookies(key ...string) CookiesProvider {
	if len(key) == 1 {
		return p.factory.Cookies(key[0])
	}
	if v := p.find(p.cookiesKeys, func(key string) interface{} {
		return p.factory.Cookies(key)
	}); v != nil {
		return v.(CookiesProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterHTTPServer(key string, httpServer HTTPServerProvider) {
	p.addKey(&p.httpServerKeys, key)
	p.factory.RegisterHTTPServer(key, httpServer)
}

func (p *AutoProvider) HTTPServer(key ...string) HTTPServerProvider {
	if len(key) == 1 {
		return p.factory.HTTPServer(key[0])
	}
	if v := p.find(p.httpServerKeys, func(key string) interface{} {
		return p.factory.HTTPServer(key)
	}); v != nil {
		return v.(HTTPServerProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterAPIServer(key string, apiServer APIServerProvider) {
	p.addKey(&p.apiServerKeys, key)
	p.factory.RegisterAPIServer(key, apiServer)
}

func (p *AutoProvider) APIServer(key ...string) APIServerProvider {
	if len(key) == 1 {
		return p.factory.APIServer(key[0])
	}
	if v := p.find(p.apiServerKeys, func(key string) interface{} {
		return p.factory.APIServer(key)
	}); v != nil {
		return v.(APIServerProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterApp(key string, app AppProvider) {
	p.addKey(&p.appKeys, key)
	p.factory.RegisterApp(key, app)
}

func (p *AutoProvider) App(key ...string) AppProvider {
	if len(key) == 1 {
		return p.factory.App(key[0])
	}
	if v := p.find(p.appKeys, func(key string) interface{} {
		return p.factory.App(key)
	}); v != nil {
		return v.(AppProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterCtx(key string, ctx CtxProvider) {
	p.addKey(&p.ctxKeys, key)
	p.factory.RegisterCtx(key, ctx)
}

func (p *AutoProvider) Ctx(key ...string) CtxProvider {
	if len(key) == 1 {
		return p.factory.Ctx(key[0])
	}
	if v := p.find(p.ctxKeys, func(key string) interface{} {
		return p.factory.Ctx(key)
	}); v != nil {
		return v.(CtxProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterError(key string, error ErrorProvider) {
	p.addKey(&p.errorKeys, key)
	p.factory.RegisterError(key, error)
}

func (p *AutoProvider) Error(key ...string) ErrorProvider {
	if len(key) == 1 {
		return p.factory.Error(key[0])
	}
	if v := p.find(p.errorKeys, func(key string) interface{} {
		return p.factory.Error(key)
	}); v != nil {
		return v.(ErrorProvider)
	}
	return nil
}

func (p *AutoProvider) RegisterI18n(key string, i18n I18nProvider) {
	p.addKey(&p.i18nKeys, key)
	p.factory.RegisterI18n(key, i18n)
}

func (p *AutoProvider) I18n(key ...string) I18nProvider {
	if len(key) == 1 {
		return p.factory.I18n(key[0])
	}
	if v := p.find(p.i18nKeys, func(key string) interface{} {
		return p.factory.I18n(key)
	}); v != nil {
		return v.(I18nProvider)
	}
	return nil
}

func NewAutoProvider() *AutoProvider {
	p := &AutoProvider{
		factory: NewProviderFactory(),
	}
	return p
}

type SessionProvider func() Session
type SessionStorageProvider func(ctx Ctx) (SessionStorage, error)
type CacheProvider func(ctx Ctx) (Cache, error)
type TemplateProvider func(ctx Ctx, rootDir string) (Template, error)
type ViewProvider func(w io.Writer, tpl Template) View
type DatabaseProvider func(ctx Ctx) (Database, error)
type DatabaseGroupProvider func(ctx Ctx) (DatabaseGroup, error)
type HTTPRouterProvider func(ctx Ctx) HTTPRouter
type ConfigProvider func() Config
type LoggerProvider func(ctx Ctx, prefix string) Logger

type ControllerProvider func(ctx Ctx) Controller
type CookiesProvider func(ctx Ctx) Cookies
type HTTPServerProvider func(ctx Ctx) HTTPServer
type APIServerProvider func(ctx Ctx, address string) (APIServer, error)
type AppProvider func(isDefault bool) App
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

func NewProviderFactory() *ProviderFactory {
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
