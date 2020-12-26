package gcore

var Providers = NewProvider()

type SessionProvider func(ctx Ctx) Session
type SessionStorageProvider func(ctx Ctx) (SessionStorage, error)
type CacheProvider func(ctx Ctx) (Cache, error)
type TemplateProvider func(ctx Ctx) (Template, error)
type ViewProvider func(ctx Ctx) View
type DatabaseProvider func(ctx Ctx) (Database, error)
type DatabaseGroupProvider func(ctx Ctx) (DatabaseGroup, error)
type HTTPRouterProvider func(ctx Ctx) HTTPRouter

type ProviderFactory struct {
	session        map[string]SessionProvider
	sessionStorage map[string]SessionStorageProvider
	cache          map[string]CacheProvider
	template       map[string]TemplateProvider
	view           map[string]ViewProvider
	database       map[string]DatabaseProvider
	databaseGroup  map[string]DatabaseGroupProvider
	httpRouter     map[string]HTTPRouterProvider
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
	}
}
