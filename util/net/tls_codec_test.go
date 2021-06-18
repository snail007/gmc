// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"crypto/tls"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTLSCodec_1(t *testing.T) {
	l, p, _ := ListenRandom("")
	l0 := NewEventListener(l)
	l0.AddCodecFactory(func(ctx Context) Codec {
		c0 := NewTLSServerCodec(&tls.Config{})
		c0.AddCertKey(testCert,testKEY)
		c0.RequireClientAuth()
		c0.AddClientCa(testCert)
		return c0
	})
	l0.OnAccept(func(ctx Context, c net.Conn) {
		c.Write([]byte("hello"))
	})
	l0.OnAcceptError(func(ctx Context, err error) {
		t.Log(err)
	})
	l0.Start()
	time.Sleep(time.Second)
	conn, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn0 := NewConn(conn)
	c1 := NewTLSClientCodec(&tls.Config{})
	c1.AddCertKey(demoCert, demoKEY)
	c1.SetVerifyServerName("test.com")
	c1.SkipVerify(true)
	conn0.AddCodec(c1)
	assert.NoError(t, conn0.Initialize())
	d, err := Read(conn0, 5)
	assert.NoError(t, err)
	assert.Equal(t, "hello", d)
}

var testKEY = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAtep4TfLTE04hXm8qgmfGbedDvxgo7bwpWS7seSC9H/IGg4J8
njYna/OZ+Ib0dArtW97o9V7jYgmja60ZxKxFvTvvWQDJEMwa/4JAHV05LiSeR0wu
q4gW8VeFaCNwtibNY78oNfg85mIVg3BD/NA1cqloIYh4G3bffLvXbkkNQYcins1Q
W3rswapmu6Y3VQOh/K5NsUudLWAEwStyVNF37n4zKFuMDz+gChOXp5FbcIcxh5vx
WomzQ6c30mQe9g2iyQMvfKbMmrVo247zs0dPuTw20+1nfnjyRjDY96xr4xcSd6eg
lNaE5p5zVLyXfviLr5P57h8XnLPspaXsrbPkIwIDAQABAoIBAQCFF8Jk5R9gpGzt
dk+XkO0wQ17hVH+9T0jBIv+Hr1gvIxd45+Lcraox5MvldHcs30HBUVkHDCE3/O0/
Pin4JkHvrQX0DAsO6wVlopnd4fKPu+LBLw+GF88RS4MjKaqw2bqzG4wD0FZeB6zN
uTlEoeA4v5Cb2Ahnr5Ta4WNAINo98bCYFrOgHLQWVsNoMRy+qWZmoK3vR8L/Mjqy
/5ZUDcJuJhcSQA+BlXTkGTiRjy4QuIMRK7HuZ59KGBRIyYAJVwP3dxKw9nNNcRui
8JEaMNuGoy3TLyeQgiUom9JuQwyOtzCuno2jgS8jKeQhuBNYVFSagDYJvC0ypTtz
hrNm4WKBAoGBANKPM/Y3VgbTKhT71dQ5bLdtq5jgA5jrsl7UkzVl2EYIIPgSRNoe
LHRbtmW4QRBu8gFjUcASkXqiqYQ+CLFxD6BNeyCQG8RZo3Grcyl04IhqhbuXEu+R
EEn1AI1G820nt2+rXp+2uG/fGsZmPQ+E2Pgvr5JClhP96tskyUppWghTAoGBAN0s
zA3MUxdVkt5ol2cwJPPTww+gHBL/gmmbYprW/VieoAGHXjVhyghpzGp+RnMdOcMb
paVMs9QorAsCHq+EZjjoLU2a5CbAdZZ21xCOTpKPR8RZgjAAJpJZFdcSp880UQ18
OU4U5zQJfrl9xxHUtWMG0WqK9GvpAe3w1/vey/rxAoGAWIk4ey5XcPU3u60NA3jF
+vcVcWm4eYOZ9AAEii5x2zitzEG6S9DmNmMd9fWc/jD4d5bwmAf2vg9Joj6HXz1A
KdKKlG2kD1L1w+UovmTTyOippPBoWO2xYLexbLZJwzsxCbaQSi4FrZytYIE66Zyd
svYyKBjxjCR3rX/xV+WmotsCgYAV/l9oO9pDZsINFc+AdlwmVvd9tUk1Zm0cfVQn
25sj1dpJbKGko03I2mR2booo5k4ZJcWqE1+KiGTbT2GnyH21yPjAT9fCNr86sCSg
w9XyYwca8l+s0EcFpJA0a+l+BFDPC3xTVGbNWOheH7DNCB7lcwceFiVKGciUVa/U
nwofsQKBgQCvrddwyUO4sw2yzV3Ny/Sm8EcREE406yBOdPZlctxGh7uI+d2ENoFK
n5FDiTMUA9VViMN2EYcSc4ucO0w/cEA3Oh71BUcu66ieSfPJ5WPQ6xFRcIxlZp7e
BL5v0wZyxrSVsGlXkaH8r9WW/7lMUVadLrXh01zsucEilwNW9/TmhA==
-----END RSA PRIVATE KEY-----
`)

// dns name: test.com
var testCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDSzCCAjOgAwIBAgIBATANBgkqhkiG9w0BAQsFADBGMQswCQYDVQQGEwJNTzER
MA8GA1UEChMIdGVzdC5jb20xETAPBgNVBAsTCHRlc3QuY29tMREwDwYDVQQDEwh0
ZXN0LmNvbTAgFw0yMTA2MTgwMTI5MjJaGA8xODUxMDkxMDAzMjAxNVowRjELMAkG
A1UEBhMCTU8xETAPBgNVBAoTCHRlc3QuY29tMREwDwYDVQQLEwh0ZXN0LmNvbTER
MA8GA1UEAxMIdGVzdC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQC16nhN8tMTTiFebyqCZ8Zt50O/GCjtvClZLux5IL0f8gaDgnyeNidr85n4hvR0
Cu1b3uj1XuNiCaNrrRnErEW9O+9ZAMkQzBr/gkAdXTkuJJ5HTC6riBbxV4VoI3C2
Js1jvyg1+DzmYhWDcEP80DVyqWghiHgbdt98u9duSQ1BhyKezVBbeuzBqma7pjdV
A6H8rk2xS50tYATBK3JU0XfufjMoW4wPP6AKE5enkVtwhzGHm/FaibNDpzfSZB72
DaLJAy98psyatWjbjvOzR0+5PDbT7Wd+ePJGMNj3rGvjFxJ3p6CU1oTmnnNUvJd+
+Iuvk/nuHxecs+ylpeyts+QjAgMBAAGjQjBAMA4GA1UdDwEB/wQEAwICpDAdBgNV
HSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDwYDVR0TAQH/BAUwAwEB/zANBgkq
hkiG9w0BAQsFAAOCAQEApC9nSTDw2ptLJN446ET6VSj/wHzquRXDg5b5sDk9+bGp
Srgx1Y5mAloCbeBpvj3ZyNyoaPccMXYWT1iNSsxogPGfnUPuI1X/rx+iNcs17+H5
GmM/JUi6gKp0mcfiK9OTCL9mm5pYnUtjlGVZkBW34dOAdNzxIZpYSqJLoVEozLij
iZaXnxtmU7bUMVqybN9126yQ8C56ZO/SIUOhMPerutTMlrIP8UOQ5CYMR8UyuDZh
JVBRWqSYvX+swsTSN4Yqy1UnaWvFVCtAlcThZnqxjI/XRUww3NUMCc6EgdLoSsI0
UoLzcN/OnQrCjReRS+yo+UmHZ0HT0IPDki0Eo3XPug==
-----END CERTIFICATE-----`)

var demoKEY = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAw+DwQPZpA+yC9+elVgNu4f2F3fPceL6YmAbGe7qq4cHTFkfK
PCdBxTASeoGAEbKWKTwJSxR47zhr+6m96jnvLAOIMSisxlh2jBfeItQGCJDJw0lQ
RwLsthdK71PAmjIIpNAQnjV3dFLnhTGXEJQIOgLkva8vGwu3OzbV0GtdtpVMJuhI
OilnAPAzvGFBvl5rFW4Pscc+xgYfaORyDBQL2LwpTWHmjJ3QDJbJlteDjyPxp99o
bgdWUuRbkn5g+MjvAyYr8JghIgUUlrJztttEUgJ8lEMKOod1ZIW0BxIsbaSCNy4P
pJHLPhoGJQs5p08A4gk2RxhnvfFJdU2ScVmKdQIDAQABAoIBAQCp/IlDNxRHjXbT
ALpg/LW7dTI5Pan1NyJhvG9/bK1jIbu4ODDvJvpSz7cZjUzBDwR1YF6IQ4n3wDUl
v1bK79/5iE8mqi/WKWsnhIcIHovl3xDZYsRB++3E0E39h+c7aXRK4y2ovqmdz1yQ
IEsC3hSNk3lCi8cLZ41p29qN9r7q9PylKLATLfqgsEcm40kulxJROOClf+TLtNy8
xfCqJo2VXVt4uQlrMnaz7evKe8oOAz4nQlx+dsuOAcoe0bv5q32akolDvk91UYQc
gwrFQCQpTg8WivoMQKZjxENlTWgiY3flLa1d6U6omRT5i5X/0JbgnKiD0hvTndcB
MIYLxpBBAoGBAOhTdyYrv6r3OkPBgcVzBC4efhkd9B2z08dvXgerfGFJFjvLs0nc
LVOeCggZArmvBrUcgM5IcHablf0Qf/sFp7kamZt2q0SP3hpVQK2nfuiL+c1LwMGb
t+xmc7drAzJkb9+FIrWXuTcG7KEaXDL5FyfMnSIgA5dle0YbL7UN4Cs9AoGBANfW
sd0JLGSxrnvzzd4bAjoeWCOYtef9Ky9Nv4CDgdGYgJK+lKqo6fkqraI5emumD6kC
9FuBNaOs+762lg6u/kjilHN/FeORQbL5kH4xYGc30p1iQ626smRAgTSkaLlGIFfL
ADRSum+/VtEDXwrOarw1LPL6ERRXcfosijPdJa+ZAoGANQopH4vJXEzI/oMFD4Ds
qWLIww81ljph1Rw1yWZ7JPK8orYknm4n4vknrSWYm6+7xklVlsKu+kUW/wlvTm3C
Ft5dx0JWY3a87CIefAbLUGf0hcwPm6PjX5McQ/moZy7K46rPe8nBvTBVgYo1FmYL
xUhPb2UDrOK8PAsk3x7l2LkCgYEAkaFXsxbscCh+3T18Kx84GnS87Y+tNRFZ4Pnp
e1G/9uaZ4elbL+b2r1r/etSjaBzMtjG7JD6DLaOa3GwfxVqHUjAnD+KwpzIsDRFc
T/kK3boJjo1tsrukgAYR564Cxves/O+IfMVQ6/NDJZXLu+PYmpKaeHsHqRzzV2RT
/3h4ZAkCgYEAq6y1102BGXKLF/HXasH8QGPBS5VMiIEx1YdhRwfpxXNZtuUdZlUV
2gCABQTq7IRBViTgkXTqM66UZFRTsKAlvz9MPuB6MH/7a7VWnYwbrDfMjnl1DM28
9iiBFM1VwxMhWcQyb5RGDfJingJJWeX7ekXYRic0jgVoTuydRGEgovQ=
-----END RSA PRIVATE KEY-----`)

// dns name: demo.com
var demoCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDSzCCAjOgAwIBAgIBATANBgkqhkiG9w0BAQsFADBGMQswCQYDVQQGEwJMUjER
MA8GA1UEChMIZGVtby5jb20xETAPBgNVBAsTCGRlbW8uY29tMREwDwYDVQQDEwhk
ZW1vLmNvbTAgFw0yMTA2MTgwMTMxNDBaGA8xODUxMDkxMDAzMjIzMlowRjELMAkG
A1UEBhMCTFIxETAPBgNVBAoTCGRlbW8uY29tMREwDwYDVQQLEwhkZW1vLmNvbTER
MA8GA1UEAxMIZGVtby5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQDD4PBA9mkD7IL356VWA27h/YXd89x4vpiYBsZ7uqrhwdMWR8o8J0HFMBJ6gYAR
spYpPAlLFHjvOGv7qb3qOe8sA4gxKKzGWHaMF94i1AYIkMnDSVBHAuy2F0rvU8Ca
Mgik0BCeNXd0UueFMZcQlAg6AuS9ry8bC7c7NtXQa122lUwm6Eg6KWcA8DO8YUG+
XmsVbg+xxz7GBh9o5HIMFAvYvClNYeaMndAMlsmW14OPI/Gn32huB1ZS5FuSfmD4
yO8DJivwmCEiBRSWsnO220RSAnyUQwo6h3VkhbQHEixtpII3Lg+kkcs+GgYlCzmn
TwDiCTZHGGe98Ul1TZJxWYp1AgMBAAGjQjBAMA4GA1UdDwEB/wQEAwICpDAdBgNV
HSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDwYDVR0TAQH/BAUwAwEB/zANBgkq
hkiG9w0BAQsFAAOCAQEADF2SdQjEKw1nfR2HRss/yCF6K7rBlmEv6gAhK3QAdiX+
tTCPsz67QLtAriWX8QbThS+60yI9rM5EFGsp+vDTt2q6kHgQteaMDhEhNYY6ic7g
Td2IoqGJlZ+Z/kdPjND6QXlN2IYBn6wXGAH2AcGtDSlQyxCEa07iUnvNCG45H/ta
dtMO82dmF+kYI87Y+zgjC6EQf2sWv2Jr/vZbeLkDwZbpFooMZEdXy2MxG2XNyNjb
6f+mGpjzS8zLBWTIgd8UoIo3t4VUkJsDkBnrP3fpyLBtEfcUoXATvgrOZX4Y6WOL
Rjs037Mzictr1H6w8EYHBm0ECWErH+Sn2YTxzKUZ/g==
-----END CERTIFICATE-----`)
