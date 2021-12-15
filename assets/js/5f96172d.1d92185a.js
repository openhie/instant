(self.webpackChunkinstant_hie_docs=self.webpackChunkinstant_hie_docs||[]).push([[62],{3905:function(e,t,n){"use strict";n.d(t,{Zo:function(){return u},kt:function(){return v}});var a=n(7294);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function o(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function c(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},i=Object.keys(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var l=a.createContext({}),s=function(e){var t=a.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):o(o({},t),e)),n},u=function(e){var t=s(e.components);return a.createElement(l.Provider,{value:t},e.children)},p={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},d=a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,i=e.originalType,l=e.parentName,u=c(e,["components","mdxType","originalType","parentName"]),d=s(n),v=r,m=d["".concat(l,".").concat(v)]||d[v]||p[v]||i;return n?a.createElement(m,o(o({ref:t},u),{},{components:n})):a.createElement(m,o({ref:t},u))}));function v(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var i=n.length,o=new Array(i);o[0]=d;var c={};for(var l in t)hasOwnProperty.call(t,l)&&(c[l]=t[l]);c.originalType=e,c.mdxType="string"==typeof e?e:r,o[1]=c;for(var s=2;s<i;s++)o[s]=n[s];return a.createElement.apply(null,o)}return a.createElement.apply(null,n)}d.displayName="MDXCreateElement"},3275:function(e,t,n){"use strict";n.r(t),n.d(t,{frontMatter:function(){return o},metadata:function(){return c},toc:function(){return l},default:function(){return u}});var a=n(2122),r=n(9756),i=(n(7294),n(3905)),o={id:"covid19-surveillance",title:"Covid19-Surveillance",sidebar_label:"Covid19-Surveillance",keywords:["Instant OpenHIE","Custom package","Covid19","Surveillance"],description:"Available Instant OpenHIE custom packages"},c={unversionedId:"use-case-packages/covid19-surveillance",id:"use-case-packages/covid19-surveillance",isDocsHomePage:!1,title:"Covid19-Surveillance",description:"Available Instant OpenHIE custom packages",source:"@site/docs/use-case-packages/covid19-surveillance.mdx",sourceDirName:"use-case-packages",slug:"/use-case-packages/covid19-surveillance",permalink:"/instant/docs/use-case-packages/covid19-surveillance",editUrl:"https://github.com/openhie/instant/tree/master/docs/docs/use-case-packages/covid19-surveillance.mdx",version:"current",sidebar_label:"Covid19-Surveillance",frontMatter:{id:"covid19-surveillance",title:"Covid19-Surveillance",sidebar_label:"Covid19-Surveillance",keywords:["Instant OpenHIE","Custom package","Covid19","Surveillance"],description:"Available Instant OpenHIE custom packages"},sidebar:"docs",previous:{title:"Covid19-Immunization",permalink:"/instant/docs/use-case-packages/covid19-immunization"},next:{title:"HIV Case Reporting",permalink:"/instant/docs/use-case-packages/hiv-case-reporting"}},l=[{value:"WHO Covid19-Surveillance",id:"who-covid19-surveillance",children:[]}],s={toc:l};function u(e){var t=e.components,n=(0,r.Z)(e,["components"]);return(0,i.kt)("wrapper",(0,a.Z)({},s,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("div",{className:"admonition admonition-info alert alert--info"},(0,i.kt)("div",{parentName:"div",className:"admonition-heading"},(0,i.kt)("h5",{parentName:"div"},(0,i.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,i.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"14",height:"16",viewBox:"0 0 14 16"},(0,i.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7zm1 3H6v5h2V4zm0 6H6v2h2v-2z"}))),"info")),(0,i.kt)("div",{parentName:"div",className:"admonition-content"},(0,i.kt)("p",{parentName:"div"},"The Instant OpenHIE architecture, codebase, and documentation are under active development and are subject to change. While we encourage adoption and extension of the Instant OpenHIE framework, we do not consider this ready for production use at this stage."))),(0,i.kt)("h2",{id:"who-covid19-surveillance"},"WHO Covid19-Surveillance"),(0,i.kt)("p",null,"This Instant OpenHIE package is available to demonstrate a country level implementation for tracking covid19 cases. The covid19-surveillance package allows for the submission of the following Covid19 related information:"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"Covid19 Case Report"),(0,i.kt)("li",{parentName:"ul"},"Covid19 Lab Result"),(0,i.kt)("li",{parentName:"ul"},"Covid-19 Case Outcome")),(0,i.kt)("p",null,"The goal is to allow for the collection of data from different data sources, by using an industry standard for the health information exchange and making the data available for analysis and reporting using the DHIS2 Covid19 dashboards."),(0,i.kt)("p",null,"The WHO Covid-19 Surveillance Package is available at the following location: ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/jembi/who-covid19-surveillance-package"},"https://github.com/jembi/who-covid19-surveillance-package")),(0,i.kt)("p",null,"This package is dependant on the following infrastructural packages:"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"HMIS"),(0,i.kt)("li",{parentName:"ul"},"Core")))}u.isMDXComponent=!0}}]);