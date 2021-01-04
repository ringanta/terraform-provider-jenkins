package jenkins

const getLocalUserCommand = `
import hudson.security.HudsonPrivateSecurityRealm
import hudson.security.HudsonPrivateSecurityRealm.Details
import hudson.tasks.Mailer
import groovy.json.JsonOutput

def result = [:]

def secRealm = jenkins.model.Jenkins.instance.getSecurityRealm()
if (!secRealm instanceof HudsonPrivateSecurityRealm) {
  result['error'] = true
  result['msg'] = 'Jenkins is not using local user database'
  result['data'] = [:]
  return println(JsonOutput.toJson(result))
}

user = secRealm.getUser('{{ .Username }}')
if (user != null) {
  	result['error'] = false
  	result['msg'] = ''
	result['data'] = [:]
  	result['data']['username'] = user.getId()
	result['data']['fullname'] = user.getFullName()
  	result['data']['password_hash'] = user.getProperty(Details.class).getPassword()
  	result['data']['email'] = user.getProperty(Mailer.UserProperty.class).getAddress()
  	result['data']['description'] = user.getDescription() != null ? user.getDescription() : ''
} else {
	result['error'] = false
  	result['msg'] = ''
  	result['data'] = [:]
}

return println(JsonOutput.toJson(result))
`

const createLocalUserCommand = `
import hudson.security.HudsonPrivateSecurityRealm
import hudson.tasks.Mailer
import groovy.json.JsonOutput

def result = [:]
def secRealm = jenkins.model.Jenkins.instance.getSecurityRealm()
if (!secRealm instanceof HudsonPrivateSecurityRealm) {
  result['error'] = true
  result['msg'] = 'Jenkins is not using local user database'
  result['data'] = [:]
  return println(JsonOutput.toJson(result))
}

def user = secRealm.createAccount('{{ .Username }}', '{{ .Password }}')
user.addProperty(new Mailer.UserProperty('{{ .Email }}'))
user.setFullName('{{ .Fullname }}')
user.setDescription('{{ .Description }}')
result['error'] = false
result['msg'] = 'User {{ .Username }} successfully created'
result['data'] = [:]

return println(JsonOutput.toJson(result))
`

const deleteLocalUserCommand = `
import hudson.security.HudsonPrivateSecurityRealm
import groovy.json.JsonOutput

def result = [:]
def secRealm = jenkins.model.Jenkins.instance.getSecurityRealm()
if (!secRealm instanceof HudsonPrivateSecurityRealm) {
  result['error'] = true
  result['msg'] = 'Jenkins is not using local user database'
  result['data'] = [:]
  return println(JsonOutput.toJson(result))
}

user = secRealm.getUser('{{ .Username }}')
user.delete()
result['error'] = false
result['msg'] = 'User {{ .Username }} successfully created'
result['data'] = [:]

return println(JsonOutput.toJson(result))
`

const getUserPermissionsCommand = `
import hudson.security.Permission
import groovy.json.JsonOutput

String shortName(Permission p) {
    p.id.tokenize('.')[-2..-1].join('/')
        .replace('Hudson','Overall')
        .replace('Computer', 'Agent')
        .replace('Item', 'Job')
		.replace('CredentialsProvider', 'Credentials')
		.replace('LockableResourcesManager', 'LockableResources')
}

def strategy = Jenkins.instance.getAuthorizationStrategy()
def permissions = []
def result = [error: false, msg: '', data: [:]]
result['data']['username'] = '{{ .Username }}'
result['data']['permissions'] = []
strategy.grantedPermissions.collect { permission, userList ->
	userList.collect { user ->
      if (user == '{{ .Username }}') {
        result['data']['permissions'].push(shortName(permission))
      }
    }
}
println(JsonOutput.toJson(result))
`

const createUserPermissionsCommand = `
import hudson.security.Permission
import groovy.json.JsonOutput

String shortName(Permission p) {
    p.id.tokenize('.')[-2..-1].join('/')
        .replace('Hudson','Overall')
        .replace('Computer', 'Agent')
        .replace('Item', 'Job')
		.replace('CredentialsProvider', 'Credentials')
		.replace('LockableResourcesManager', 'LockableResources')
}

def permissionIds = Permission.all.findAll { permission ->
    def nonConfigurablePerms = ['RunScripts', 'UploadPlugins', 'ConfigureUpdateCenter']
    permission.enabled &&
        !permission.id.startsWith('hudson.security.Permission') &&
        !(true in nonConfigurablePerms.collect { permission.id.endsWith(it) })
}.collect { permission ->
    [ (shortName(permission)): permission ]
}.sum()

def strategy = Jenkins.instance.getAuthorizationStrategy()
def result = [error: false, msg: '', data: [:]]
def user_permissions = [{{range .Permissions}}'{{.}}',{{end}}]

user_permissions.collect {
	strategy.add(permissionIds[it], '{{ .Username }}')
}
Jenkins.instance.save()
result['msg'] = 'Permissions for user {{ .Username }} is created'

println(JsonOutput.toJson(result))
`

const updateUserPermissionsCommand = `
import hudson.security.Permission
import groovy.json.JsonOutput

String shortName(Permission p) {
    p.id.tokenize('.')[-2..-1].join('/')
        .replace('Hudson','Overall')
        .replace('Computer', 'Agent')
        .replace('Item', 'Job')
		.replace('CredentialsProvider', 'Credentials')
		.replace('LockableResourcesManager', 'LockableResources')
}

def permissionIds = Permission.all.findAll { permission ->
    def nonConfigurablePerms = ['RunScripts', 'UploadPlugins', 'ConfigureUpdateCenter']
    permission.enabled &&
        !permission.id.startsWith('hudson.security.Permission') &&
        !(true in nonConfigurablePerms.collect { permission.id.endsWith(it) })
}.collect { permission ->
    [ (shortName(permission)): permission ]
}.sum()

def strategy = Jenkins.instance.getAuthorizationStrategy()
def result = [error: false, msg: '', data: [:]]
def user_permissions = [{{range .Permissions}}'{{.}}',{{end}}]
user_permissions.removeAll([null])

user_permissions.collect {
	strategy.add(permissionIds[it], '{{ .Username }}')
}

strategy.grantedPermissions.collect { permission, userList ->
	if (!user_permissions.contains(shortName(permission))) {
		userList.remove('{{ .Username }}')
	}
}

Jenkins.instance.save()
result['msg'] = 'Permissions of user {{ .Username }} is updated'

println(JsonOutput.toJson(result))
`

const deleteUserPermissionsCommand = `
import groovy.json.JsonOutput

def strategy = Jenkins.instance.getAuthorizationStrategy()
def result = [error: false, msg: '', data: [:]]
strategy.grantedPermissions.collect { permission, userList ->
	userList.remove('{{ .Username }}')
}
Jenkins.instance.save()
result['msg'] = 'User {{ .Username }} has been removed from the global matrix authorization'

println(JsonOutput.toJson(result))
`
