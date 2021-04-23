// Code generated by statik. DO NOT EDIT.

package embed

import (
	"github.com/rakyll/statik/fs"
)


func init() {
	data := "PK\x03\x04\x14\x00\x08\x00\x00\x00\xa9\x9b\xb4P\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00client.js.tmplUT\x05\x00\x01\x0e\x85\xc5^{{define \"client\"}}\n{{- if .Services}}\n//\n// Client\n//\n{{ range .Services}}\n{{exportKeyword}}class {{.Name}} {\n  constructor(hostname, fetch) {\n    this.path = '/rpc/{{.Name}}/'\n    this.hostname = hostname\n    this.fetch = fetch\n  }\n\n  url(name) {\n    return this.hostname + this.path + name\n  }\n  {{range .Methods}}\n  {{.Name | methodName}} = ({{.Inputs | methodInputs}}) => {\n    return this.fetch(\n      this.url('{{.Name}}'),\n      {{- if .Inputs | len}}\n      createHTTPRequest(args, headers)\n      {{- else}}\n      createHTTPRequest({}, headers)\n      {{- end}}\n    ).then((res) => {\n      return buildResponse(res).then(_data => {\n        return {\n        {{- $outputsCount := .Outputs|len -}}\n        {{- range $i, $output := .Outputs}}\n          {{$output | newOutputArgResponse}}{{listComma $i $outputsCount}}\n        {{- end}}\n        }\n      })\n    })\n  }\n  {{end}}\n}\n{{end -}}\n{{end -}}\n{{end}}PK\x07\x08&\xd2\xd7\x13\x8c\x03\x00\x00\x8c\x03\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\xa9\x9b\xb4P\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x16\x00	\x00client_helpers.js.tmplUT\x05\x00\x01\x0e\x85\xc5^{{define \"client_helpers\"}}\nconst createHTTPRequest = (body = {}, headers = {}) => {\n  return {\n    method: 'POST',\n    headers: { ...headers, 'Content-Type': 'application/json' },\n    body: JSON.stringify(body || {})\n  }\n}\n\nconst buildResponse = (res) => {\n  return res.text().then(text => {\n    let data\n    try {\n      data = JSON.parse(text)\n    } catch(err) {\n      throw { code: 'unknown', msg: `expecting JSON, got: ${text}`, status: res.status }\n    }\n    if (!res.ok) {\n      throw data // webrpc error response\n    }\n    return data\n  })\n}\n{{end}}\nPK\x07\x08\xb2\x9b\x81\xf5.\x02\x00\x00.\x02\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00K\x8e\x97R\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x11\x00	\x00proto.gen.js.tmplUT\x05\x00\x01\xdf\x08\x83`{{- define \"proto\" -}}\n// {{.Name}} {{.SchemaVersion}} {{.SchemaHash}}\n\n// --\n// This file has been generated by https://github.com/webrpc/webrpc using gen/javascript\n// Do not edit by hand. Update your webrpc schema and re-generate.\n\n// WebRPC description and code-gen version\nexport const WebRPCVersion = \"{{.WebRPCVersion}}\"\n\n// Schema version of your RIDL schema\nexport const WebRPCSchemaVersion = \"{{.SchemaVersion}}\"\n\n{{template \"types\" .}}\n{{- if .TargetOpts.Client}}\n  {{template \"client\" .}}\n  {{template \"client_helpers\" .}}\n{{- end}}\n{{- if .TargetOpts.Server}}\n  {{template \"server\" .}}\n{{- end}}\n{{- end}}\nPK\x07\x08\xa9.&\xf3k\x02\x00\x00k\x02\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\xa9\x9b\xb4P\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00server.js.tmplUT\x05\x00\x01\x0e\x85\xc5^{{define \"server\"}}\n\n{{- if .Services}}\n//\n// Server\n//\n\nclass WebRPCError extends Error {\n    constructor(msg = \"error\", statusCode) {\n        super(\"webrpc eror: \" + msg);\n\n        this.statusCode = statusCode\n    }\n}\n\nimport express from 'express'\n\n    {{- range .Services}}\n        {{$name := .Name}}\n        {{$serviceName := .Name | serviceInterfaceName}}\n\n        export const create{{$serviceName}}App = (serviceImplementation) => {\n            const app = express();\n\n            app.use(express.json())\n\n            app.post('/*', async (req, res) => {\n                const requestPath = req.baseUrl + req.path\n\n                if (!req.body) {\n                    res.status(400).send(\"webrpc error: missing body\");\n\n                    return\n                }\n\n                switch(requestPath) {\n                    {{range .Methods}}\n\n                    case \"/rpc/{{$name}}/{{.Name}}\": {                        \n                        try {\n                            {{ range .Inputs }}\n                                {{- if not .Optional}}\n                                    if (!(\"{{ .Name }}\" in req.body)) {\n                                        throw new WebRPCError(\"Missing Argument `{{ .Name }}`\")\n                                    }\n                                {{end -}}\n\n                                if (typeof req.body[\"{{.Name}}\"] !== \"{{ .Type | jsFieldType }}\") {\n                                    throw new WebRPCError(\"Invalid arg: {{ .Name }}, got type \" + typeof req.body[\"{{ .Name }}\"] + \" expected \" + \"{{ .Type | jsFieldType }}\", 400);\n                                }\n                            {{end}}\n\n                            const response = await serviceImplementation[\"{{.Name}}\"](req.body);\n\n                            {{ range .Outputs}}\n                                if (!(\"{{ .Name }}\" in response)) {\n                                    throw new WebRPCError(\"internal\", 500);\n                                }\n                            {{end}}\n\n                            res.status(200).json(response);\n                        } catch (err) {\n                            if (err instanceof WebRPCError) {\n                                const statusCode = err.statusCode || 400\n                                const message = err.message\n\n                                res.status(statusCode).json({\n                                    msg: message,\n                                    status: statusCode,\n                                    code: \"\"\n                                });\n\n                                return\n                            }\n\n                            if (err.message) {\n                                res.status(400).send(err.message);\n\n                                return;\n                            }\n\n                            res.status(400).end();\n                        }\n                    }\n\n                    return;\n                    {{end}}\n\n                    default: {\n                        res.status(404).end()\n                    }\n                }\n            });\n\n            return app;\n        };\n    {{- end}}\n{{end -}}\n{{end}}\nPK\x07\x08>E\\\xb6s\x0c\x00\x00s\x0c\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\xa9\x9b\xb4P\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0d\x00	\x00types.js.tmplUT\x05\x00\x01\x0e\x85\xc5^{{define \"types\"}}\n//\n// Types\n//\n{{ if .Messages -}}\n{{range .Messages -}}\n\n{{if .Type | isEnum -}}\n{{$enumName := .Name}}\n{{exportKeyword}}var {{$enumName}};\n(function ({{$enumName}}) {\n{{- range $i, $field := .Fields}}\n  {{$enumName}}[\"{{$field.Name}}\"] = \"{{$field.Name}}\"\n{{- end}}\n})({{$enumName}} || ({{$enumName}} = {}))\n{{end -}}\n\n{{- if .Type | isStruct  }}\n{{exportKeyword}}class {{.Name}} {\n  constructor(_data) {\n    this._data = {}\n    if (_data) {\n      {{range .Fields -}}\n      this._data['{{. | exportedJSONField}}'] = _data['{{. | exportedJSONField}}']\n      {{end}}\n    }\n  }\n  {{ range .Fields -}}\n  get {{. | exportedJSONField}}() {\n    return this._data['{{. | exportedJSONField }}']\n  }\n  set {{. | exportedJSONField}}(value) {\n    this._data['{{. | exportedJSONField}}'] = value\n  }\n  {{end}}\n  toJSON() {\n    return this._data\n  }\n}\n{{end -}}\n{{end -}}\n{{end -}}\n\n{{end}}\nPK\x07\x08r\x06\xac\x87\x82\x03\x00\x00\x82\x03\x00\x00PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xa9\x9b\xb4P&\xd2\xd7\x13\x8c\x03\x00\x00\x8c\x03\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\x80\x81\x00\x00\x00\x00client.js.tmplUT\x05\x00\x01\x0e\x85\xc5^PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xa9\x9b\xb4P\xb2\x9b\x81\xf5.\x02\x00\x00.\x02\x00\x00\x16\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\x80\x81\xd1\x03\x00\x00client_helpers.js.tmplUT\x05\x00\x01\x0e\x85\xc5^PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00K\x8e\x97R\xa9.&\xf3k\x02\x00\x00k\x02\x00\x00\x11\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81L\x06\x00\x00proto.gen.js.tmplUT\x05\x00\x01\xdf\x08\x83`PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xa9\x9b\xb4P>E\\\xb6s\x0c\x00\x00s\x0c\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\x80\x81\xff\x08\x00\x00server.js.tmplUT\x05\x00\x01\x0e\x85\xc5^PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xa9\x9b\xb4Pr\x06\xac\x87\x82\x03\x00\x00\x82\x03\x00\x00\x0d\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\x80\x81\xb7\x15\x00\x00types.js.tmplUT\x05\x00\x01\x0e\x85\xc5^PK\x05\x06\x00\x00\x00\x00\x05\x00\x05\x00c\x01\x00\x00}\x19\x00\x00\x00\x00"
		fs.Register(data)
	}
	