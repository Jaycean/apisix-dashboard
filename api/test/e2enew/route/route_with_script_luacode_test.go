/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package route

import (
	"net/http"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/gomega"

	"github.com/apisix/manager-api/test/e2enew/base"
)

var t = ginkgo.GinkgoT()

var _ = ginkgo.Describe("route with script lucacode", func() {
	ginkgo.It("clean APISIX error log", func() {
		base.CleanAPISIXErrorLog()
	})
	table.DescribeTable("test route with script lucacode",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		table.Entry("create route with script of valid lua code", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello",
				"upstream": {
						"type": "roundrobin",
						"nodes": {
								"` + base.UpstreamIp + `:1980": 1
						}
				},
				"script": "local _M = {} \n function _M.access(api_ctx) \n ngx.log(ngx.WARN,\"hit access phase\") \n end \nreturn _M"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		table.Entry("get the route", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hit access phase",
			Sleep:        base.SleepTime,
		}),
		table.Entry("hit the route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world\n",
		}),
		table.Entry("update route with script of valid lua code", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello",
				"upstream": {
						"type": "roundrobin",
						"nodes": {
								"` + base.UpstreamIp + `:1981": 1
						}
				},
				"script": "local _M = {} \n function _M.access(api_ctx) \n ngx.log(ngx.WARN,\"hit access phase\") \n end \nreturn _M"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		table.Entry("update route with script of invalid lua code", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello",
				"upstream": {
						"type": "roundrobin",
						"nodes": {
								"` + base.UpstreamIp + `:1980": 1
						}
				},
				"script": "local _M = {} \n function _M.access(api_ctx) \n ngx.log(ngx.WARN,\"hit access phase\")"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		table.Entry("delete the route (r1)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		}),
		table.Entry("hit the route just delete", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		}),
		table.Entry("create route with script of invalid lua code", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello",
				"upstream": {
						"type": "roundrobin",
						"nodes": {
								"` + base.UpstreamIp + `:1980": 1
						}
				},
				"script": "local _M = {} \n function _M.access(api_ctx) \n ngx.log(ngx.WARN,\"hit access phase\")"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
	)
	ginkgo.It("verify the log generated by script set in Step-3 above", func() {
		//sleep for process log
		time.Sleep(1500 * time.Millisecond)

		//verify the log generated by script set in Step-3 above
		logContent := base.ReadAPISIXErrorLog()
		gomega.Expect(logContent).Should(gomega.ContainSubstring(`\"hit access phase\"`))

		//clean log
		base.CleanAPISIXErrorLog()
	})
})

var _ = ginkgo.Describe("route with script id", func() {
	ginkgo.It("clean APISIX error log", func() {
		base.CleanAPISIXErrorLog()
	})
	table.DescribeTable("test route with script id",
		func(tc base.HttpTestCase) {
			base.RunTestCase(tc)
		},
		table.Entry("create route with invalid script_id - not equal to id", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/routes",
			Body: `{
				"id": "r1",
				"name": "route1",
				"uri": "/hello",
				"upstream": {
						"type": "roundrobin",
						"nodes": {
								"` + base.UpstreamIp + `:1980": 1
						}
				},
				"script": "local _M = {} \n function _M.access(api_ctx) \n ngx.log(ngx.WARN,\"hit access phase\") \n end \nreturn _M",
				"script_id": "not-r1"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		table.Entry("create route with invalid script_id - set script_id but without id", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/routes",
			Body: `{
				"name": "route1",
				"uri": "/hello",
				"upstream": {
						"type": "roundrobin",
						"nodes": {
								"` + base.UpstreamIp + `:1980": 1
						}
				},
				"script": "local _M = {} \n function _M.access(api_ctx) \n ngx.log(ngx.WARN,\"hit access phase\") \n end \nreturn _M",
				"script_id": "r1"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		table.Entry("create route with invalid script_id - set script_id but without script", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/routes",
			Body: `{
				"id": "r1",
				"name": "route1",
				"uri": "/hello",
				"upstream": {
						"type": "roundrobin",
						"nodes": {
								"` + base.UpstreamIp + `:1980": 1
						}
				},
				"script_id": "r1"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		table.Entry("create route with valid script_id", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPost,
			Path:   "/apisix/admin/routes",
			Body: `{
					"id": "r1",
					"name": "route1",
					"uri": "/hello",
					"upstream": {
						"type": "roundrobin",
						"nodes": {
							"` + base.UpstreamIp + `:1980": 1
						}
					},
					"script": "local _M = {} \n function _M.access(api_ctx) \n ngx.log(ngx.WARN,\"hit access phase\") \n end \nreturn _M",
					"script_id": "r1"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		table.Entry("get the route", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodGet,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "\"script_id\":\"r1\"",
			Sleep:        base.SleepTime,
		}),
		table.Entry("hit the route", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello world\n",
		}),
		table.Entry("update route with valid script_id", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello",
				"upstream": {
						"type": "roundrobin",
						"nodes": {
								"` + base.UpstreamIp + `:1981": 1
						}
				},
				"script": "local _M = {} \n function _M.access(api_ctx) \n ngx.log(ngx.WARN,\"hit access phase\") \n end \nreturn _M",
				"script_id": "r1"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
		}),
		table.Entry("update route with invalid script_id - not equal to id", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello",
				"upstream": {
						"type": "roundrobin",
						"nodes": {
								"` + base.UpstreamIp + `:1980": 1
						}
				},
				"script": "local _M = {} \n function _M.access(api_ctx) \n ngx.log(ngx.WARN,\"hit access phase\")",
				"script_id": "not-r1"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		table.Entry("update route with invalid script_id - set script_id but without script", base.HttpTestCase{
			Object: base.ManagerApiExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/r1",
			Body: `{
				"name": "route1",
				"uri": "/hello",
				"upstream": {
					"type": "roundrobin",
					"nodes": {
							"` + base.UpstreamIp + `:1980": 1
					}
				},
				"script_id": "r1"
			}`,
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusBadRequest,
		}),
		table.Entry("delete the route (r1)", base.HttpTestCase{
			Object:       base.ManagerApiExpect(),
			Method:       http.MethodDelete,
			Path:         "/apisix/admin/routes/r1",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusOK,
			Sleep:        base.SleepTime,
		}),
		table.Entry("hit the route just delete", base.HttpTestCase{
			Object:       base.APISIXExpect(),
			Method:       http.MethodGet,
			Path:         "/hello",
			Headers:      map[string]string{"Authorization": base.GetToken()},
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"error_msg":"404 Route Not Found"}`,
			Sleep:        base.SleepTime,
		}),
	)
	ginkgo.It("verify the log generated by script set in Step-4 above", func() {
		//sleep for process log
		time.Sleep(1500 * time.Millisecond)

		//verify the log generated by script set in Step-4 above
		logContent := base.ReadAPISIXErrorLog()
		gomega.Expect(logContent).Should(gomega.ContainSubstring(`\"hit access phase\"`))

		//clean log
		base.CleanAPISIXErrorLog()
	})
})
