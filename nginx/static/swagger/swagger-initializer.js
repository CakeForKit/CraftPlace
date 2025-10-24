
window.onload = function() {
  // Build a system
  const ui = SwaggerUIBundle({
    url: "http://localhost:8080/swagger/doc.json",
    dom_id: '#swagger-ui',
    validatorUrl: null,
    oauth2RedirectUrl: `${window.location.protocol}//${window.location.host}${window.location.pathname.split('/').slice(0, window.location.pathname.split('/').length - 1).join('/')}/oauth2-redirect.html`,
    persistAuthorization: false,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
	layout: "StandaloneLayout",
    docExpansion: "list",
	deepLinking: true,
	defaultModelsExpandDepth: 1
  })

  const defaultClientId = "";
  if (defaultClientId) {
    ui.initOAuth({
      clientId: defaultClientId,
      usePkceWithAuthorizationCodeGrant: false
    })
  }

  window.ui = ui
}
