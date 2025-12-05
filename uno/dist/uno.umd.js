(function(y,w){typeof exports=="object"&&typeof module<"u"?w(require("alpinejs"),require("htmx.org")):typeof define=="function"&&define.amd?define(["alpinejs","htmx.org"],w):(y=typeof globalThis<"u"?globalThis:y||self,w(y.Alpine,y.htmx))})(this,(function(y,w){"use strict";const j={box:{ok:"OK",cancel:"Cancel",loading:"Loading",close:"Close"},valid:{required:"Required",pattern:"Invalid format",min:"At least {n} characters",validname:"Letters, numbers, - and _ only",en:"English letters only",domain:"Invalid domain format",mobile:"Invalid mobile number",int:"Integer required",email:"Invalid email address",date:"Invalid date format (YYYY-MM-DD)",url:"Invalid URL",pwd:"At least 8 chars, 1 letter, 1 number",strongpwd:"At least 8 chars, 1 letter, 1 number, 1 special char"},charCount:"Remaining: {count}"},z={box:{ok:"确认",cancel:"取消",loading:"加载中",close:"关闭"},valid:{required:"必填",pattern:"格式无效",min:"至少 {n} 个字符",validname:"仅限字母、数字、- 和 _",en:"仅限英文字母",domain:"无效的域名格式",mobile:"无效的手机号",int:"需要整数",email:"无效的邮箱地址",date:"无效的日期格式 (YYYY-MM-DD)",url:"无效的 URL",pwd:"至少 8 位，包含 1 个字母和 1 个数字",strongpwd:"至少 8 位，包含字母、数字和特殊字符"},charCount:"剩余: {count}"},B={box:{ok:"確認",cancel:"取消",loading:"加載中",close:"關閉"},valid:{required:"必填",pattern:"格式無效",min:"至少 {n} 個字符",validname:"僅限字母、數字、- 和 _",en:"僅限英文字母",domain:"無效的域名格式",mobile:"無效的手機號",int:"需要整數",email:"無效的電子郵件地址",date:"無效的日期格式 (YYYY-MM-DD)",url:"無效的 URL",pwd:"至少 8 位，包含 1 個字母和 1 個數字",strongpwd:"至少 8 位，包含字母、數字和特殊字符"},charCount:"剩餘: {count}"};const T={en:j,"zh-CN":z,"zh-TW":B},I=Object.keys(T);function O(){if(typeof localStorage<"u"){const s=localStorage.getItem("uno.lang");if(s&&I.includes(s))return s}if(typeof navigator<"u"){const s=navigator.language,d=s.startsWith("zh")?s.includes("TW")||s.includes("Hant")?"zh-TW":"zh-CN":"en";return I.includes(d)?d:"en"}return"en"}function R(s,d){const o=d.split(".");let t=s;for(const n of o){if(!t||typeof t!="object"||!(n in t))return;t=t[n]}return t}function M(s,d){for(const o in d)d[o]&&typeof d[o]=="object"?(s[o]||Object.assign(s,{[o]:{}}),M(s[o],d[o])):Object.assign(s,{[o]:d[o]})}function q(s){s.store("i18n",{locale:O(),messages:T,init(){typeof document<"u"&&(document.documentElement.lang=this.locale)},t(d,o={}){const t=this.locale,n=this.messages[t]||this.messages.en,e=R(n,d);return(typeof e=="string"?e:d).replace(/\{(\w+)\}/g,(a,r)=>r in o?String(o[r]):`{${r}}`)},setLang(d){if(!this.messages[d]){console.warn(`[Uno] Language ${d} not found`);return}this.locale=d,typeof localStorage<"u"&&localStorage.setItem("uno.lang",d),typeof document<"u"&&(document.documentElement.lang=d)},extend(d){for(const o in d)this.messages[o]?M(this.messages[o],d[o]):this.messages[o]=d[o]}})}const p={info(...s){console.log("[Uno Info]",...s)},warn(...s){console.warn("[Uno Warn]",...s)},error(...s){console.error("[Uno Error]",...s)}};function S(s,d,o){const t=s.getBoundingClientRect(),n=d.getBoundingClientRect(),e=n.width,l=n.height,a=t.width,r=t.height,i={top:t.top,left:t.left},c=window.innerWidth,f=window.innerHeight,u=0,m=0,g=f,h=c,v=(L,F)=>{if(F){if(L+e>h)return h-e-10;if(L<m)return m+10}else{if(L+l>g)return g-l-10;if(L<u)return u+10}return L};let x,b,C=o;switch(o){case"left":b=i.left-e-5,b<m&&(b=i.left+a+5,C="right"),x=v(i.top+r/2-l/2,!1);break;case"top":x=i.top-l-8,x<u&&(x=i.top+r+8,C="bottom"),b=i.left+(a-e)/2,b=v(b,!0);break;case"right":b=i.left+a+5,b+e>h&&(b=i.left-e-5,C="left"),x=v(i.top+r/2-l/2,!1);break;case"bottom":default:x=i.top+r+8,x+l>g&&(x=i.top-l-8,C="top"),b=i.left+(a-e)/2,b=v(b,!0);break}return{top:x,left:b,position:C}}const U=s=>{if(!s){p.error("Tooltip: Alpine is undefined. Plugin initialization failed.");return}s.store("tooltip",{visible:!1,content:"",placement:"top",top:-9999,left:-9999,theme:"dark",triggerEl:null,async show(o,t,n={}){this.triggerEl=o,this.content=t,this.placement=n.placement||"top",this.theme=n.theme||"dark",this.visible=!0,this.ensureContainer(),await s.nextTick(),this.updatePosition()},hide(){this.visible=!1,this.triggerEl=null},updatePosition(){if(!this.triggerEl)return;const o=document.getElementById("uno-tooltip-content");if(!o)return;const t=S(this.triggerEl,o,this.placement);this.top=t.top,this.left=t.left},ensureContainer(){document.getElementById("uno-tooltip-container")||d()}});const d=()=>{if(!document.getElementById("uno-tooltip-container"))try{const o=document.createElement("div");o.id="uno-tooltip-container",o.setAttribute("x-data","{}"),o.innerHTML=`
          <div id="uno-tooltip-content"
               x-show="$store.tooltip.visible"
               x-transition:enter="transition ease-out duration-200"
               x-transition:enter-start="opacity-0 scale-95"
               x-transition:enter-end="opacity-100 scale-100"
               x-transition:leave="transition ease-in duration-150"
               x-transition:leave-start="opacity-100 scale-100"
               x-transition:leave-end="opacity-0 scale-95"
               class="fixed z-[9999] px-3 py-2 text-sm font-medium rounded-lg shadow-sm whitespace-nowrap pointer-events-none border"
               :class="$store.tooltip.theme === 'light' ? 'bg-white text-gray-900 border-gray-200' : 'bg-slate-800 text-white border-slate-700'"
               :style="\`top: \${$store.tooltip.top}px; left: \${$store.tooltip.left}px\`"
               style="display: none;">
               <span x-html="$store.tooltip.content"></span>
          </div>
        `,document.body?document.body.appendChild(o):document.documentElement.appendChild(o)}catch(o){p.error("Tooltip: Failed to inject container",o)}};document.body?d():document.addEventListener("DOMContentLoaded",d),s.directive("tooltip",(o,{expression:t,modifiers:n},{evaluateLater:e,effect:l,cleanup:a})=>{const r=e(t);let i={content:"",placement:n[0]||"top",trigger:n.includes("click")?"click":"hover",theme:"dark"};n.includes("top")&&(i.placement="top"),n.includes("bottom")&&(i.placement="bottom"),n.includes("left")&&(i.placement="left"),n.includes("right")&&(i.placement="right"),n.includes("light")&&(i.theme="light");const c=g=>{if(typeof g=="string"?i.content=g:typeof g=="object"&&(i={...i,...g}),i.content&&i.content.startsWith("#")&&!i.content.includes(" ")){const h=document.querySelector(i.content);h&&(i.content=h.innerHTML)}else if(i.target){const h=document.querySelector(i.target);h&&(i.content=h.innerHTML)}};l(()=>{r(g=>{c(g)})});const f=()=>s.store("tooltip").show(o,i.content,i),u=()=>s.store("tooltip").hide(),m=()=>{const g=s.store("tooltip");g.visible&&g.triggerEl===o?g.hide():f()};if(i.trigger==="click"){o.addEventListener("click",m);const g=h=>{const v=s.store("tooltip");v.visible&&v.triggerEl===o&&!o.contains(h.target)&&v.hide()};window.addEventListener("click",g),a(()=>window.removeEventListener("click",g))}else o.addEventListener("mouseenter",f),o.addEventListener("mouseleave",u),o.addEventListener("focus",f),o.addEventListener("blur",u);a(()=>{i.trigger==="click"?o.removeEventListener("click",m):(o.removeEventListener("mouseenter",f),o.removeEventListener("mouseleave",u),o.removeEventListener("focus",f),o.removeEventListener("blur",u))})})};const H=s=>{s.store("notify",{items:[],nextId:1,activeRequests:0,activeProgress:0,show(t,n={}){const{type:e="info",timeout:l=3e3,position:a="top-right",closeIndicator:r=!1,showCloseButton:i=!0}=n,c=this.nextId++,f={id:c,text:t,type:e,position:a,closeIndicator:r,showCloseButton:i,visible:!0};this.items.unshift(f),l>0&&setTimeout(()=>{this.close(c)},l),this.ensureContainer()},close(t){this.items.find(e=>e.id===t)&&(this.items=this.items.filter(e=>e.id!==t))},success(t,n={}){const e=typeof n=="number"?{timeout:n}:n;this.show(t,{type:"success",...e})},warning(t,n={}){const e=typeof n=="number"?{timeout:n}:n;this.show(t,{type:"warning",...e})},error(t,n={}){const e=typeof n=="number"?{timeout:n}:n;this.show(t,{type:"error",...e})},info(t,n={}){const e=typeof n=="number"?{timeout:n}:n;this.show(t,{type:"info",...e})},loading(t,n={}){t?this.activeRequests++:this.activeRequests=Math.max(0,this.activeRequests-1),this.ensureContainer()},progress(t){t?this.activeProgress++:this.activeProgress=Math.max(0,this.activeProgress-1),this.ensureContainer()},ensureContainer(){document.getElementById("uno-notify-container")||d()}});const d=()=>{try{const t=document.createElement("div");t.id="uno-notify-container",t.setAttribute("x-data","{}");const n=`
        <style>
          @keyframes uno-progress-indeterminate {
            0% { left: -35%; width: 35%; }
            100% { left: 100%; width: 35%; }
          }
          .uno-progress-anim {
            animation: uno-progress-indeterminate 1.5s infinite linear;
          }
        </style>

        <!-- Progress Bar (Top Fixed) -->
        <div id="uno-progress-bar"
             x-show="$store.notify.activeProgress > 0"
             class="fixed top-0 left-0 w-full h-1 bg-blue-100 z-[200]"
             style="display: none;">
             <div class="absolute top-0 h-full bg-blue-600 uno-progress-anim"></div>
        </div>

        <!-- Loading Indicator (Top Right Fixed) -->
        <div id="uno-loading-indicator" 
             x-show="$store.notify.activeRequests > 0"
             class="fixed top-4 right-4 z-[150] bg-blue-100 border-blue-500 text-blue-700 p-2 rounded-full shadow-md flex items-center dark:bg-blue-900 dark:text-blue-100"
             style="display: none;">
             <svg class="animate-spin -ml-1 mr-2 h-5 w-5 text-blue-500 dark:text-blue-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
             <span>Loading</span>
        </div>

        <!-- Top Right -->
        <div class="fixed top-0 right-0 z-[100] p-4 flex flex-col space-y-2 items-end pointer-events-none">
            <template x-for="item in $store.notify.items.filter(i => i.position === 'top-right')" :key="item.id">
                ${o()}
            </template>
        </div>

        <!-- Top Left -->
        <div class="fixed top-0 left-0 z-[100] p-4 flex flex-col space-y-2 items-start pointer-events-none">
            <template x-for="item in $store.notify.items.filter(i => i.position === 'top-left')" :key="item.id">
                ${o()}
            </template>
        </div>

        <!-- Bottom Right -->
        <div class="fixed bottom-0 right-0 z-[100] p-4 flex flex-col-reverse space-y-2 space-y-reverse items-end pointer-events-none">
            <template x-for="item in $store.notify.items.filter(i => i.position === 'bottom-right')" :key="item.id">
                ${o()}
            </template>
        </div>

        <!-- Bottom Left -->
        <div class="fixed bottom-0 left-0 z-[100] p-4 flex flex-col-reverse space-y-2 space-y-reverse items-start pointer-events-none">
            <template x-for="item in $store.notify.items.filter(i => i.position === 'bottom-left')" :key="item.id">
                ${o()}
            </template>
        </div>
      `;t.innerHTML=n,document.body.appendChild(t)}catch(t){p.error("Notify: Failed to inject container",t)}},o=()=>`
    <div x-show="true"
         data-testid="notify-item"
         x-transition:enter="transition ease-out duration-300"
         x-transition:enter-start="opacity-0 transform scale-90 translate-y-2"
         x-transition:enter-end="opacity-100 transform scale-100 translate-y-0"
         x-transition:leave="transition ease-in duration-300"
         x-transition:leave-start="opacity-100 transform scale-100 translate-y-0"
         x-transition:leave-end="opacity-0 transform scale-90 translate-y-2"
         class="pointer-events-auto relative w-80 max-w-full border-l-4 pl-3 pr-3 py-4 rounded shadow-sm dark:shadow-md overflow-hidden"
         :class="{
            'bg-green-100 dark:bg-green-800 border-green-500 text-green-700 dark:text-green-100': item.type === 'success',
            'bg-yellow-100 dark:bg-yellow-800 border-yellow-500 text-yellow-700 dark:text-yellow-100': item.type === 'warning',
            'bg-red-100 dark:bg-red-800 border-red-500 text-red-700 dark:text-red-100': item.type === 'error',
            'bg-blue-100 dark:bg-blue-800 border-blue-500 text-blue-700 dark:text-blue-100': item.type === 'info'
         }"
         role="alert">
         
         <div class="flex items-start justify-between">
            <div x-text="item.text" class="text-sm font-medium"></div>
            <button x-show="item.showCloseButton" 
                    @click="$store.notify.close(item.id)" 
                    class="ml-4 opacity-50 hover:opacity-100 transition-opacity focus:outline-none">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
            </button>
         </div>
    </div>
    `;document.body?d():document.addEventListener("DOMContentLoaded",d)};let k=0,$="";function D(){k===0&&($=document.body.style.overflow||"",document.body.style.overflow="hidden"),k++}function P(){k=Math.max(0,k-1),k===0&&(document.body.style.overflow=$||"",$="")}async function _(s){const d=await fetch(s);return(d.headers.get("content-type")||"").includes("application/json")?d.json():d.text()}const N=s=>{s.store("box",{items:[],init(){s.effect(()=>{this.items.filter(t=>t.modal).length>0?D():P()})},async open(o,t={}){let n="dialog",e=t;typeof o=="object"?(e=o,n=e.type||"dialog"):n=o;const l=e.boxId||Math.random().toString(36).slice(2,10),a={id:l,type:n,title:e.title||"",content:e.content||"",width:e.width||400,height:e.height,modal:e.modal!==!1,overlayOpacity:e.overlayOpacity??.5,position:e.position||"center",showCloseButton:!0,okText:e.okButtonText,cancelText:e.cancelButtonText,onOk:e.onOk,onCancel:e.onCancel,showOk:n==="confirm"||n==="message"||n==="dialog"&&e.showOkButton,showCancel:n==="confirm"||n==="dialog"&&e.showCancelButton,loading:!1,visible:!0},r=s.store("i18n");if(a.okText||(a.okText=r?r.t("box.ok"):"OK"),a.cancelText||(a.cancelText=r?r.t("box.cancel"):"Cancel"),e.targetId){const i=document.getElementById(e.targetId);i?a.content=i.innerHTML:(p.warn("Box: Target element not found:",e.targetId),a.content='<div class="p-4 text-red-500">Error: Content target not found</div>')}else if(e.ajaxUrl){a.loading=!0,this.add(a);try{const i=await _(e.ajaxUrl),c=typeof i=="object"?`<pre class="bg-gray-100 p-3 rounded overflow-x-auto whitespace-pre-wrap">${JSON.stringify(i,null,2)}</pre>`:String(i),f=this.items.find(u=>u.id===l);f&&(f.content=c,f.loading=!1)}catch(i){p.error("Box: Ajax failed",i);const c=this.items.find(f=>f.id===l);c&&(c.content='<div class="text-red-500">Error loading content</div>',c.loading=!1)}return l}else e.iframeUrl&&(a.content=`<iframe frameborder="0" class="w-full h-full" src="${e.iframeUrl}"></iframe>`);return this.add(a),l},add(o){this.items.push(o),this.ensureContainer()},close(o){this.items=this.items.filter(t=>t.id!==o)},closeAll(){this.items=[]},ensureContainer(){document.getElementById("uno-box-container")||d()}});const d=()=>{try{const o=document.createElement("div");o.id="uno-box-container",o.setAttribute("x-data","{}");const t=`
        <template x-for="(item, index) in $store.box.items" :key="item.id">
            <div :id="'uno-box-outer-container' + item.id"
                 class="fixed inset-0 flex items-center justify-center transition-opacity duration-300"
                 :style="\`z-index: \${100 + index}\`"
                 role="dialog"
                 aria-modal="true">
                 
                 <!-- Overlay: Covers viewport, blocks interaction, does NOT close modal on click -->
                 <div :id="'uno-box-overlay' + item.id"
                      class="absolute inset-0 bg-black transition-opacity duration-300"
                      :style="\`opacity: \${item.overlayOpacity}\`"
                      x-show="item.modal"
                      x-transition:enter="ease-out duration-300"
                      x-transition:enter-start="opacity-0"
                      x-transition:enter-end="opacity-100"
                      x-transition:leave="ease-in duration-200"
                      x-transition:leave-start="opacity-100"
                      x-transition:leave-end="opacity-0"></div>

                 <!-- Box Container -->
                 <div :id="'uno-box-container' + item.id"
                      class="relative bg-white dark:bg-gray-900 rounded-lg shadow-2xl ring-1 ring-black/5 overflow-hidden flex flex-col max-h-[90vh] transition-all duration-300 transform"
                      :style="\`width: \${item.width}px; \${item.height ? 'height:' + item.height + 'px' : ''}\`"
                      x-show="item.visible"
                      x-transition:enter="ease-out duration-300"
                      x-transition:enter-start="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                      x-transition:enter-end="opacity-100 translate-y-0 sm:scale-100"
                      x-transition:leave="ease-in duration-200"
                      x-transition:leave-start="opacity-100 translate-y-0 sm:scale-100"
                      x-transition:leave-end="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95">
                      
                      <!-- Header -->
                      <div :id="'uno-box-header' + item.id" x-show="item.title" class="flex justify-between items-center p-4 border-b border-slate-200 dark:border-slate-700">
                          <h3 :id="'uno-box-title' + item.id" class="text-lg font-semibold text-slate-800 dark:text-slate-200" x-text="item.title"></h3>
                          <button :id="'uno-box-close' + item.id" @click="$store.box.close(item.id)" class="text-slate-400 hover:text-slate-500 focus:outline-none">
                              <span class="sr-only">Close</span>
                              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                              </svg>
                          </button>
                      </div>
                      
                      <!-- Close button if no title -->
                      <button x-show="!item.title && item.showCloseButton" 
                              :id="'uno-box-close' + item.id"
                              @click="$store.box.close(item.id)" 
                              class="absolute top-2 right-2 z-10 text-slate-400 hover:text-slate-500 focus:outline-none">
                          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                          </svg>
                      </button>

                      <!-- Content -->
                      <div :id="'uno-box-wrapper' + item.id" class="p-4 overflow-auto flex-1 text-slate-600 dark:text-slate-300">
                          <div x-show="item.loading" class="flex justify-center items-center h-full">
                              <svg class="animate-spin h-8 w-8 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                              </svg>
                          </div>
                          <div :id="'uno-box-content' + item.id" x-show="!item.loading" x-html="item.content"></div>
                      </div>

                      <!-- Footer -->
                      <div :id="'uno-box-footer' + item.id" x-show="item.showOk || item.showCancel" class="p-4 border-t border-slate-200 dark:border-slate-700 flex justify-end space-x-3 bg-slate-50 dark:bg-slate-800">
                          <button :id="'uno-box-cancel' + item.id"
                                  x-show="item.showCancel"
                                  @click="item.onCancel ? (item.onCancel(), $store.box.close(item.id)) : $store.box.close(item.id)"
                                  x-text="item.cancelText"
                                  class="px-4 py-2 bg-slate-100 dark:bg-slate-700 border border-slate-200 dark:border-slate-600 rounded-md text-sm font-medium text-slate-800 dark:text-slate-200 hover:bg-slate-200 dark:hover:bg-slate-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors">
                          </button>
                          <button :id="'uno-box-ok' + item.id"
                                  x-show="item.showOk"
                                  @click="() => {
                                      const shouldClose = item.onOk ? item.onOk() : true;
                                      if (shouldClose !== false) $store.box.close(item.id);
                                  }"
                                  x-text="item.okText"
                                  class="px-4 py-2 bg-blue-600 border border-transparent rounded-md text-sm font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                          </button>
                      </div>
                 </div>
            </div>
        </template>
      `;o.innerHTML=t,document.body.appendChild(o)}catch(o){p.error("Box: Failed to inject container",o)}};document.body?d():document.addEventListener("DOMContentLoaded",d),document.addEventListener("click",o=>{const t=o.target.closest("[data-box-type]");if(!t)return;const n=t.getAttribute("data-box-type");if(!["dialog","confirm","message"].includes(n))return;o.preventDefault();const e=t.dataset,l={title:e.boxTitle,content:e.boxContent,width:e.boxWidth?parseInt(e.boxWidth):void 0,height:e.boxHeight?parseInt(e.boxHeight):void 0,ajaxUrl:e.boxAjaxUrl,iframeUrl:e.boxIframeUrl,targetId:e.boxTargetId,position:e.boxPosition,okButtonText:e.boxOkText,cancelButtonText:e.boxCancelText,overlayOpacity:e.boxOverlayOpacity?parseFloat(e.boxOverlayOpacity):void 0,showCloseButton:e.boxShowClose!=="false",boxId:e.boxId,cache:e.boxCache==="true"};e.boxModal==="false"&&(l.modal=!1),s.store("box").open(n,l)})},A={validname:/^[a-zA-Z0-9_-]+$/,en:/^[a-zA-Z]+$/,domain:new RegExp("^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\\.)+[A-Za-z]{2,6}$"),mobile:/^\+?[\d\s-]{7,}$/,int:/^-?\d+$/,email:/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/,date:/^(\d{4}-\d{2}-\d{2}|\d{2}\/\d{2}\/\d{4}|\d{2}\.\d{2}\.\d{4})$/,url:/^(https?:\/\/)?([\da-z.-]+)\.([a-z.-]{2,6})([/\w .-]*)*\/?$/i,pwd:/^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/,strongpwd:/^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/};function W(s){const d=s.match(/^\/(.*?)\/([gimsuy]*)$/);return d?{pattern:d[1],flags:d[2]}:{pattern:s,flags:""}}function Y(s){try{return s?JSON.parse(s):[]}catch{return[]}}const Z=s=>{const d={invalid:`url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 12 12' width='12' height='12' fill='none' stroke='%23ef4444'%3e%3ccircle cx='6' cy='6' r='4.5'/%3e%3cpath stroke-linecap='round' d='M5.8 3.6h.4L6 6.5zM6 8.2af.6.6 0 110-1.2.6.6 0 010 1.2z'/%3e%3c/svg%3e")`,valid:`url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 8'%3e%3cpath fill='%2322c55e' d='M2.3 6.73L.6 4.53c-.4-1.04.46-1.4 1.1-.8l1.1 1.4 3.4-3.8c.6-.63 1.6-.27 1.2.7l-4 4.6c-.43.5-.8.4-1.1.1z'/%3e%3c/svg%3e")`},o=(t,n)=>{const e=s.store("i18n");if(!e)return t;const l={required:"Required",pattern:"Invalid format",min:`At least ${n?.n} characters`,max:`Max ${n?.n} characters`,int:"Integer required",email:"Invalid email address",url:"Invalid URL"};return e.t(`valid.${t}`,n)||l[t]||t};s.store("valid",{callbacks:{},registerCallback(t,n){typeof n=="function"&&(this.callbacks[t]=n)}}),s.data("validForm",()=>({fields:{},init(){this.$el.noValidate=!0},register(t,n,e=[]){if(this.fields[t])return;const l=[...e];if((n.hasAttribute("required")||n.dataset.validRequired==="true")&&l.push({type:"required"}),(n.type==="email"||n.dataset.validType==="email")&&l.push("email"),(n.type==="url"||n.dataset.validType==="url")&&l.push("url"),n.dataset.validType){const a=n.dataset.validType;a!=="email"&&a!=="url"&&l.push(a)}n.dataset.validPattern&&l.push({type:"regexp",expression:n.dataset.validPattern}),this.fields[t]={el:n,rules:l,value:n.value,error:null,valid:!0,dirty:!1,touched:!1,validating:!1}},update(t,n){this.fields[t]&&(this.fields[t].value=n,this.fields[t].dirty=!0,this.validate(t))},touch(t){this.fields[t]&&(this.fields[t].touched=!0,this.validate(t))},async validate(t){const n=this.fields[t];if(!n)return!0;n.validating=!0,n.error=null;const e=(n.value||"").trim();let l=null;for(const a of n.rules){const r=a.type||a,i=typeof a=="object"?a:{};if(!e){if(r==="required"){l=i.msg||o("required");break}if(r==="emptyquit")break;continue}if(A[r])A[r].test(e)||(l=i.msg||o(r)||o("pattern"));else if(r==="regexp"&&i.expression){const{pattern:c,flags:f}=W(i.expression);new RegExp(c,f||"i").test(e)||(l=i.msg||o("pattern"))}else if(r==="limit")i.min&&e.length<i.min&&(l=i.msg||o("min",{n:i.min})),i.max&&e.length>i.max&&(l=i.msg||o("max",{n:i.max}));else if(r==="range"){const c=parseFloat(e);i.min!==void 0&&c<i.min&&(l=i.msg||`Min ${i.min}`),i.max!==void 0&&c>i.max&&(l=i.msg||`Max ${i.max}`)}else if(r==="compareto"){const c=document.getElementById(i.destid);c&&e!==c.value&&(l=i.msg||"Not match")}else if(r==="tag"){const c=i.mark||",",f=e.split(c).filter(u=>u.trim());i.max&&f.length>i.max&&(l=i.msg||`Max ${i.max} items`)}if(l)break;if(r==="callback"&&i.fn){const c=s.store("valid").callbacks[i.fn];if(c)try{const f=await c(n.el.id,e);f===!1||typeof f=="object"&&f.valid===!1?l=(typeof f=="object"?f.message:"")||"Invalid":typeof f=="string"&&(l=f)}catch{}}if(l)break}return n.error=l,n.valid=!l,n.validating=!1,n.valid},async validateAll(){const t=Object.keys(this.fields);return(await Promise.all(t.map(e=>this.validate(e)))).every(e=>e)},async submit(t){if(t&&(t.preventDefault(),t.stopPropagation()),Object.keys(this.fields).forEach(e=>{this.fields[e].touched=!0}),await this.validateAll())this.$el.dispatchEvent(new CustomEvent("valid-submit",{bubbles:!0})),!this.$el.hasAttribute("hx-post")&&!this.$el.hasAttribute("hx-get")&&this.$el.submit();else{const e=Object.keys(this.fields).find(l=>!this.fields[l].valid);e&&this.fields[e].el&&this.fields[e].el.focus()}}})),s.directive("validate",(t,n,{effect:e,cleanup:l})=>{const a=t.getAttribute("name");if(!a)return;let r=[];const i=t.getAttribute("valid");i&&(r=Y(i)),e(()=>{const u=s.$data(t);u&&typeof u.register=="function"&&u.register(a,t,r)});const c=()=>{const u=s.$data(t);if(u&&u.update&&u.update(a,t.value),r.find(m=>(m.type==="limit"||m==="limit")&&m.show===!0)){const m=r.find(g=>g.type==="limit"||g==="limit");m&&m.max&&t.value.length>m.max&&(t.value=t.value.substring(0,m.max),u&&u.update&&u.update(a,t.value))}},f=()=>{const u=s.$data(t);u&&u.touch&&u.touch(a)};t.addEventListener("input",c),t.addEventListener("blur",f),e(()=>{const u=s.$data(t);if(u&&u.fields&&u.fields[a]){const m=u.fields[a],g=!!m.error,h=m.valid&&m.dirty&&!m.error;t.classList.toggle("is-invalid",g),t.classList.toggle("is-valid",h),t.classList.toggle("validating",m.validating),g?(t.classList.add("!border-red-500","pr-10","bg-no-repeat","bg-[right_0.75rem_center]","bg-[length:1rem_1rem]"),t.classList.remove("!border-green-500"),t.style.backgroundImage=d.invalid):h?(t.classList.add("!border-green-500","pr-10","bg-no-repeat","bg-[right_0.75rem_center]","bg-[length:1rem_1rem]"),t.classList.remove("!border-red-500"),t.style.backgroundImage=d.valid):(t.classList.remove("!border-red-500","!border-green-500","pr-10","bg-no-repeat","bg-[right_0.75rem_center]","bg-[length:1rem_1rem]"),t.style.backgroundImage="");const v=t.nextElementSibling;if(v&&v.tagName==="DIV"){const x=v.querySelector(".animate-spin")||t.parentNode.querySelector(".animate-spin");x&&(x.style.display=m.validating?"block":"")}}}),l(()=>{t.removeEventListener("input",c),t.removeEventListener("blur",f)})}),s.directive("valid-error",(t,{expression:n},{effect:e})=>{const l=n;t.style.display="none",e(()=>{const a=s.$data(t);if(a&&a.fields&&a.fields[l]){const i=a.fields[l].error;i?(t.textContent=i,t.style.display="block",t.classList.add("text-red-600","text-xs"),t.classList.remove("text-gray-400")):t.style.display="none"}})}),s.directive("valid-count",(t,{expression:n},{effect:e})=>{const l=n;e(()=>{const a=s.$data(t);if(a&&a.fields&&a.fields[l]){const r=a.fields[l],i=r.rules.find(c=>c.type==="limit"||c==="limit");if(i&&i.max){const c=(r.value||"").length,f=i.max-c,u=s.store("i18n");t.textContent=u&&u.t("charCount",{count:f})||`Remaining: ${f}`,t.style.display="block"}else t.style.display="none"}})})};const V=s=>{const d={activeClass:"text-blue-600 dark:text-blue-400 border-b-2 border-blue-500 dark:border-blue-400",inactiveClass:"text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200",baseTabClass:"px-4 py-2 transition-colors duration-150 cursor-pointer"},o=(a,r)=>{const i=d.activeClass.split(" "),c=d.inactiveClass.split(" ");r?(a.classList.add("current"),a.classList.remove(...c),a.classList.add(...i)):(a.classList.remove("current"),a.classList.remove(...i),a.classList.add(...c))},t=(a,r)=>{const i=document.querySelectorAll(`[data-tab-group="${a}"], [tabgroup="${a}"]`);i.forEach(f=>{const m=(f.getAttribute("data-tab-target")||f.getAttribute("tabtarget"))===r;o(f,m)});const c=new Set;i.forEach(f=>{const u=f.getAttribute("data-tab-target")||f.getAttribute("tabtarget");u&&c.add(u)}),c.forEach(f=>{const u=document.getElementById(f);u&&(f===r?(u.classList.remove("hidden"),u.style.display=""):u.classList.add("hidden"))})},n=()=>{const a=window.location.hash.substring(1);if(!a)return;const r=document.querySelector(`[data-tab-hash="${a}"], [rel="${a}"]`);if(r){const i=r.getAttribute("data-tab-group")||r.getAttribute("tabgroup"),c=r.getAttribute("data-tab-target")||r.getAttribute("tabtarget");i&&c&&t(i,c)}},e=()=>{document.querySelectorAll("[data-tab-group], [tabgroup]").forEach(r=>{if(!r.className.includes("transition-colors")){const c=d.baseTabClass.split(" ");r.classList.add(...c)}const i=r.classList.contains("current")||r.getAttribute("data-tab-active")==="true";o(r,i)})};document.addEventListener("click",a=>{const r=a.target.closest("[data-tab-target], [tabtarget]");if(!r)return;const i=r.getAttribute("data-tab-group")||r.getAttribute("tabgroup"),c=r.getAttribute("data-tab-target")||r.getAttribute("tabtarget"),f=r.getAttribute("data-tab-hash")||r.getAttribute("rel"),u=r.getAttribute("data-tab-url");if(i&&c){if(a.preventDefault(),u){const m=document.getElementById(c);m&&!m.getAttribute("data-tab-loaded")&&(m.innerHTML=`<div class="animate-pulse p-4 space-y-2">
             <div class="h-4 bg-gray-200 rounded w-3/4"></div>
             <div class="h-4 bg-gray-200 rounded w-1/2"></div>
           </div>`,fetch(u).then(g=>{if(!g.ok)throw new Error("Network response was not ok");return g.text()}).then(g=>{m.innerHTML=g,m.setAttribute("data-tab-loaded","true")}).catch(g=>{m.innerHTML=`<div class="text-red-500 p-4">Error loading content: ${g.message}</div>`}))}t(i,c),f&&history.replaceState(null,null,"#"+f)}});const l=()=>{n(),e(),window.addEventListener("hashchange",n)};document.readyState==="loading"?document.addEventListener("DOMContentLoaded",l):l(),s.store("tabs",{switch(a,r){t(a,r)}})};window.Alpine=y,window.htmx=w,p.info("Uno: Script loaded. Import Alpine:",!!y,"Window Alpine:",!!window.Alpine);const E=s=>{if(!s){p.error("Uno: Alpine is undefined, cannot register plugins");return}if(s.store&&s.store("tooltip")){p.info("Uno: Plugins already registered");return}p.info("Uno: Registering plugins"),q(s),U(s),H(s),N(s),Z(s),V(s)};(()=>{let s=y||window.Alpine;s?(p.info("Uno: Alpine found, registering"),E(s)):(p.info("Uno: Alpine not found, waiting for alpine:init"),document.addEventListener("alpine:init",()=>{p.info("Uno: alpine:init fired"),E(window.Alpine)}),document.addEventListener("DOMContentLoaded",()=>{window.Alpine&&(p.info("Uno: DOMContentLoaded fallback"),E(window.Alpine))}))})()}));
