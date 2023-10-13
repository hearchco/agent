package swisscows

import (
	"github.com/robertkrimen/otto"
	"github.com/rs/zerolog/log"
)

func performMagic(input string) string {
	vm := otto.New()
	if err := vm.Set("input_string", input); err != nil {
		log.Error().Err(err).Msg("swisscows: failed setting otto vm")
	}
	if _, err := vm.Run(`
		function i(_){var n,r,e,$="",g=-1;if(_&&_.length)for(e=_.length;(g+=1)<e;)n=_.charCodeAt(g),r=g+1<e?_.charCodeAt(g+1):0,55296<=n&&n<=56319&&56320<=r&&r<=57343&&(n=65536+((1023&n)<<10)+(1023&r),g+=1),n<=127?$+=String.fromCharCode(n):n<=2047?$+=String.fromCharCode(192|n>>>6&31,128|63&n):n<=65535?$+=String.fromCharCode(224|n>>>12&15,128|n>>>6&63,128|63&n):n<=2097151&&($+=String.fromCharCode(240|n>>>18&7,128|n>>>12&63,128|n>>>6&63,128|63&n));return $}function a(_,n){var r=(65535&_)+(65535&n);return(_>>16)+(n>>16)+(r>>16)<<16|65535&r}function s(_,n){return _<<n|_>>>32-n}function u(_,n){for(var r,e=n?"0123456789ABCDEF":"0123456789abcdef",$="",g=0,C=_.length;g<C;g+=1)r=_.charCodeAt(g),$+=e.charAt(r>>>4&15)+e.charAt(15&r);return $}function c(_){var n,r=32*_.length,e="";for(n=0;n<r;n+=8)e+=String.fromCharCode(_[n>>5]>>>24-n%32&255);return e}function l(_){var n,r=32*_.length,e="";for(n=0;n<r;n+=8)e+=String.fromCharCode(_[n>>5]>>>n%32&255);return e}function f(_){var n,r=8*_.length,e=Array(_.length>>2),$=e.length;for(n=0;n<$;n+=1)e[n]=0;for(n=0;n<r;n+=8)e[n>>5]|=(255&_.charCodeAt(n/8))<<n%32;return e}function p(_){var n,r=8*_.length,e=Array(_.length>>2),$=e.length;for(n=0;n<$;n+=1)e[n]=0;for(n=0;n<r;n+=8)e[n>>5]|=(255&_.charCodeAt(n/8))<<24-n%32;return e}function d(_,n){var r,e,$,g,C,A,v,m,b=n.length,y=[];for(g=(A=Array(Math.ceil(_.length/2))).length,r=0;r<g;r+=1)A[r]=_.charCodeAt(2*r)<<8|_.charCodeAt(2*r+1);for(;A.length>0;){for(C=[],$=0,r=0;r<A.length;r+=1)$=($<<16)+A[r],$-=(e=Math.floor($/b))*b,(C.length>0||e>0)&&(C[C.length]=e);y[y.length]=$,A=C}for(v="",r=y.length-1;r>=0;r--)v+=n.charAt(y[r]);for(m=Math.ceil(8*_.length/(Math.log(n.length)/Math.log(2))),r=v.length;r<m;r+=1)v=n[0]+v;return v}function h(_,n){var r,e,$,g="",C=_.length;for(n=n||"=",r=0;r<C;r+=3)for($=_.charCodeAt(r)<<16|(r+1<C?_.charCodeAt(r+1)<<8:0)|(r+2<C?_.charCodeAt(r+2):0),e=0;e<4;e+=1)8*r+6*e>8*_.length?g+=n:g+="ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/".charAt($>>>6*(3-e)&63);return g}var o={SHA256:function(_){_&&"boolean"==typeof _.uppercase&&_.uppercase;var n,r=_&&"string"==typeof _.pad?_.pad:"=",e=!_||"boolean"!=typeof _.utf8||_.utf8;function $(_,n){return c(x(p(_=n?i(_):_),8*_.length))}function g(_,n){_=e?i(_):_,n=e?i(n):n;var r,$=0,g=p(_),C=Array(16),A=Array(16);for(g.length>16&&(g=x(g,8*_.length));$<16;$+=1)C[$]=909522486^g[$],A[$]=1549556828^g[$];return r=x(C.concat(p(n)),512+8*n.length),c(x(A.concat(r),768))}function C(_,n){return _>>>n|_<<32-n}function A(_,n){return _>>>n}function v(_,n,r){return _&n^~_&r}function m(_,n,r){return _&n^_&r^n&r}function b(_){return C(_,2)^C(_,13)^C(_,22)}function y(_){return C(_,6)^C(_,11)^C(_,25)}function w(_){var n;return C(_,7)^C(_,18)^(n=_)>>>3}function x(_,r){var e,$,g,A,v,x,E,S,q,B,D,F,H,K,T,U,G,J=[1779033703,-1150833019,1013904242,-1521486534,1359893119,-1694144372,528734635,1541459225],L=Array(64);for(_[r>>5]|=128<<24-r%32,_[15+(r+64>>9<<4)]=r,H=0;H<_.length;H+=16){for(v=J[0],x=J[1],E=J[2],S=J[3],q=J[4],B=J[5],D=J[6],F=J[7],K=0;K<64;K+=1){L[K]=K<16?_[K+H]:a(a(a(C(G=L[K-2],17)^C(G,19)^(A=G)>>>10,L[K-7]),w(L[K-15])),L[K-16]),T=a(a(a(a(F,y(q)),(e=q,$=B,e&$^~e&(g=D))),n[K]),L[K]),U=a(b(v),m(v,x,E)),F=D,D=B,B=q,q=a(S,T),S=E,E=x,x=v,v=a(T,U)}J[0]=a(v,J[0]),J[1]=a(x,J[1]),J[2]=a(E,J[2]),J[3]=a(S,J[3]),J[4]=a(q,J[4]),J[5]=a(B,J[5]),J[6]=a(D,J[6]),J[7]=a(F,J[7])}return J}this.hex=function(_){return u($(_,e))},this.b64=function(_){return h($(_,e),r)},this.any=function(_,n){return d($(_,e),n)},this.raw=function(_){return $(_,e)},this.hex_hmac=function(_,n){return u(g(_,n))},this.b64_hmac=function(_,n){return h(g(_,n),r)},this.any_hmac=function(_,n,r){return d(g(_,n),r)},this.vm_test=function(){return"900150983cd24fb0d6963f7d28e17f72"===hex("abc").toLowerCase()},this.setUpperCase=function(_){return this},this.setPad=function(_){return r=_||r,this},this.setUTF8=function(_){return"boolean"==typeof _&&(e=_),this},n=[1116352408,1899447441,-1245643825,-373957723,961987163,1508970993,-1841331548,-1424204075,-670586216,310598401,607225278,1426881987,1925078388,-2132889090,-1680079193,-1046744716,-459576895,-272742522,264347078,604807628,770255983,1249150122,1555081692,1996064986,-1740746414,-1473132947,-1341970488,-1084653625,-958395405,-710438585,113926993,338241895,666307205,773529912,1294757372,1396182291,1695183700,1986661051,-2117940946,-1838011259,-1564481375,-1474664885,-1035236496,-949202525,-778901479,-694614492,-200395387,275423344,430227734,506948616,659060556,883997877,958139571,1322822218,1537002063,1747873779,1955562222,2024104815,-2067236844,-1933114872,-1866530822,-1538233109,-1090935817,-965641998]}},t=new o.SHA256,res=t.b64(input_string);
	`); err != nil {
		log.Error().Err(err).Msg("swisscows: failed running otto vm")
	}

	var result string = ""
	if value, err := vm.Get("res"); err == nil {
		if value_str, err := value.ToString(); err == nil {
			result = value_str
		} else {
			log.Error().Msgf("Swisscows: MagicPerformer - couldn't convert returned value to string. Error: %v", err)
		}
	} else {
		log.Error().Msgf("Swisscows: MagicPerformer - couldn't get result from javascript VM. Error: %v", err)
	}

	return result
}
