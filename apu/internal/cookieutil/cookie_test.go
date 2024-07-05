package cookieutil

import (
	"fmt"
	"testing"
)

func TestCookieUtils_Str2Map(t *testing.T) {
	s := "rewardsn=; wxtokenkey=777; wxuin=2237150243; devicetype=iMacMacBookPro171OSXOSX14.5build(23F79); version=13080512; lang=en; appmsg_token=1277_oiLzzwBRwCGhJjH4eKK5i7tGM7ZeBLhTjAPe1X7yMc_m-zRfdrf2clen3moyDG0x-Fs303gLNyhfRgHB; pass_ticket=EhK9vQxP7Cr+iKS3OuDhqLi0qMVOygmlFiXjMpHYk3Jzzlwg4xRFty+Z3asWD6wM; wap_sid2=CKPo4KoIEooBeV9IQVRXRGRPTWdOS1V3MGM2NTVNTy1vUndZVC1oUDRRNE5RdHdGYjA1RFNCNk9UTU5TbXdjZk1ua0JycGswTm4tQlZoclN0Y1lUSTNTOENSOXp3c0NId2dQejhpbldfaVZ0clJxYzJXWlNVV1dpVTUxNlNlMTg5Rkhud3BCX1lIU3o5b1NBQUF+MNWHmLQGOA1AAQ=="
	m := StrToMap(s)
	fmt.Println(m)
}
