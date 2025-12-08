(function(y,w){typeof exports=="object"&&typeof module<"u"?w(require("alpinejs"),require("htmx.org")):typeof define=="function"&&define.amd?define(["alpinejs","htmx.org"],w):(y=typeof globalThis<"u"?globalThis:y||self,w(y.Alpine,y.htmx))})(this,(function(y,w){"use strict";const B={box:{ok:"OK",cancel:"Cancel",loading:"Loading",close:"Close"},valid:{required:"Required",pattern:"Invalid format",min:"At least {n} characters",validname:"Letters, numbers, - and _ only",en:"English letters only",domain:"Invalid domain format",mobile:"Invalid mobile number",int:"Integer required",email:"Invalid email address",date:"Invalid date format (YYYY-MM-DD)",url:"Invalid URL",pwd:"At least 8 chars, 1 letter, 1 number",strongpwd:"At least 8 chars, 1 letter, 1 number, 1 special char"},charCount:"Remaining: {count}"},O={box:{ok:"确认",cancel:"取消",loading:"加载中",close:"关闭"},valid:{required:"必填",pattern:"格式无效",min:"至少 {n} 个字符",validname:"仅限字母、数字、- 和 _",en:"仅限英文字母",domain:"无效的域名格式",mobile:"无效的手机号",int:"需要整数",email:"无效的邮箱地址",date:"无效的日期格式 (YYYY-MM-DD)",url:"无效的 URL",pwd:"至少 8 位，包含 1 个字母和 1 个数字",strongpwd:"至少 8 位，包含字母、数字和特殊字符"},charCount:"剩余: {count}"},R={box:{ok:"確認",cancel:"取消",loading:"加載中",close:"關閉"},valid:{required:"必填",pattern:"格式無效",min:"至少 {n} 個字符",validname:"僅限字母、數字、- 和 _",en:"僅限英文字母",domain:"無效的域名格式",mobile:"無效的手機號",int:"需要整數",email:"無效的電子郵件地址",date:"無效的日期格式 (YYYY-MM-DD)",url:"無效的 URL",pwd:"至少 8 位，包含 1 個字母和 1 個數字",strongpwd:"至少 8 位，包含字母、數字和特殊字符"},charCount:"剩餘: {count}"};const I={en:B,"zh-CN":O,"zh-TW":R},j=Object.keys(I);function q(){if(typeof localStorage<"u"){const s=localStorage.getItem("uno.lang");if(s&&j.includes(s))return s}if(typeof navigator<"u"){const s=navigator.language,a=s.startsWith("zh")?s.includes("TW")||s.includes("Hant")?"zh-TW":"zh-CN":"en";return j.includes(a)?a:"en"}return"en"}function A(s,a){const n=a.split(".");let t=s;for(const i of n){if(!t||typeof t!="object"||!(i in t))return;t=t[i]}return t}function M(s,a){for(const n in a)a[n]&&typeof a[n]=="object"?(s[n]||Object.assign(s,{[n]:{}}),M(s[n],a[n])):Object.assign(s,{[n]:a[n]})}function S(s){s.store("i18n",{locale:q(),messages:I,init(){typeof document<"u"&&(document.documentElement.lang=this.locale)},t(a,n={}){const t=this.locale,i=this.messages[t]||this.messages.en,e=A(i,a);return(typeof e=="string"?e:a).replace(/\{(\w+)\}/g,(l,d)=>d in n?String(n[d]):`{${d}}`)},setLang(a){if(!this.messages[a]){console.warn(`[Uno] Language ${a} not found`);return}this.locale=a,typeof localStorage<"u"&&localStorage.setItem("uno.lang",a),typeof document<"u"&&(document.documentElement.lang=a)},extend(a){for(const n in a)this.messages[n]?M(this.messages[n],a[n]):this.messages[n]=a[n]}})}const x={info(...s){console.log("[Uno Info]",...s)},warn(...s){console.warn("[Uno Warn]",...s)},error(...s){console.error("[Uno Error]",...s)}};function H(s,a,n){const t=s.getBoundingClientRect(),i=a.getBoundingClientRect(),e=i.width,r=i.height,l=t.width,d=t.height,o={top:t.top,left:t.left},u=window.innerWidth,f=window.innerHeight,c=0,h=0,m=f,p=u,v=(C,J)=>{if(J){if(C+e>p)return p-e-10;if(C<h)return h+10}else{if(C+r>m)return m-r-10;if(C<c)return c+10}return C};let b,g,$=n;switch(n){case"left":g=o.left-e-5,g<h&&(g=o.left+l+5,$="right"),b=v(o.top+d/2-r/2,!1);break;case"top":b=o.top-r-8,b<c&&(b=o.top+d+8,$="bottom"),g=o.left+(l-e)/2,g=v(g,!0);break;case"right":g=o.left+l+5,g+e>p&&(g=o.left-e-5,$="left"),b=v(o.top+d/2-r/2,!1);break;case"bottom":default:b=o.top+d+8,b+r>m&&(b=o.top-r-8,$="top"),g=o.left+(l-e)/2,g=v(g,!0);break}return{top:b,left:g,position:$}}const P=s=>{if(!s){x.error("Tooltip: Alpine is undefined. Plugin initialization failed.");return}s.store("tooltip",{visible:!1,content:"",placement:"top",top:-9999,left:-9999,theme:"dark",triggerEl:null,async show(n,t,i={}){this.triggerEl=n,this.content=t,this.placement=i.placement||"top",this.theme=i.theme||"dark",this.visible=!0,this.ensureContainer(),await s.nextTick(),this.updatePosition()},hide(){this.visible=!1,this.triggerEl=null},updatePosition(){if(!this.triggerEl)return;const n=document.getElementById("uno-tooltip-content");if(!n)return;const t=H(this.triggerEl,n,this.placement);this.top=t.top,this.left=t.left},ensureContainer(){document.getElementById("uno-tooltip-container")||a()}});const a=()=>{if(!document.getElementById("uno-tooltip-container"))try{const n=document.createElement("div");n.id="uno-tooltip-container",n.setAttribute("x-data","{}"),n.innerHTML=`
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
        `,document.body?document.body.appendChild(n):document.documentElement.appendChild(n)}catch(n){x.error("Tooltip: Failed to inject container",n)}};document.body?a():document.addEventListener("DOMContentLoaded",a),s.directive("tooltip",(n,{expression:t,modifiers:i},{evaluateLater:e,effect:r,cleanup:l})=>{const d=e(t);let o={content:"",placement:i[0]||"top",trigger:i.includes("click")?"click":"hover",theme:"dark"};i.includes("top")&&(o.placement="top"),i.includes("bottom")&&(o.placement="bottom"),i.includes("left")&&(o.placement="left"),i.includes("right")&&(o.placement="right"),i.includes("light")&&(o.theme="light");const u=m=>{if(typeof m=="string"?o.content=m:typeof m=="object"&&(o={...o,...m}),o.content&&o.content.startsWith("#")&&!o.content.includes(" ")){const p=document.querySelector(o.content);p&&(o.content=p.innerHTML)}else if(o.target){const p=document.querySelector(o.target);p&&(o.content=p.innerHTML)}};r(()=>{d(m=>{u(m)})});const f=()=>s.store("tooltip").show(n,o.content,o),c=()=>s.store("tooltip").hide(),h=()=>{const m=s.store("tooltip");m.visible&&m.triggerEl===n?m.hide():f()};if(o.trigger==="click"){n.addEventListener("click",h);const m=p=>{const v=s.store("tooltip");v.visible&&v.triggerEl===n&&!n.contains(p.target)&&v.hide()};window.addEventListener("click",m),l(()=>window.removeEventListener("click",m))}else n.addEventListener("mouseenter",f),n.addEventListener("mouseleave",c),n.addEventListener("focus",f),n.addEventListener("blur",c);l(()=>{o.trigger==="click"?n.removeEventListener("click",h):(n.removeEventListener("mouseenter",f),n.removeEventListener("mouseleave",c),n.removeEventListener("focus",f),n.removeEventListener("blur",c))})})};const U=s=>{s.store("notify",{items:[],nextId:1,activeRequests:0,activeProgress:0,show(t,i={}){const{type:e="info",timeout:r=3e3,position:l="top-right",closeIndicator:d=!1,showCloseButton:o=!0}=i,u=this.nextId++,f={id:u,text:t,type:e,position:l,closeIndicator:d,showCloseButton:o,visible:!0};this.items.unshift(f),r>0&&setTimeout(()=>{this.close(u)},r),this.ensureContainer()},close(t){this.items.find(e=>e.id===t)&&(this.items=this.items.filter(e=>e.id!==t))},success(t,i={}){const e=typeof i=="number"?{timeout:i}:i;this.show(t,{type:"success",...e})},warning(t,i={}){const e=typeof i=="number"?{timeout:i}:i;this.show(t,{type:"warning",...e})},error(t,i={}){const e=typeof i=="number"?{timeout:i}:i;this.show(t,{type:"error",...e})},info(t,i={}){const e=typeof i=="number"?{timeout:i}:i;this.show(t,{type:"info",...e})},loading(t,i={}){t?this.activeRequests++:this.activeRequests=Math.max(0,this.activeRequests-1),this.ensureContainer()},progress(t){t?this.activeProgress++:this.activeProgress=Math.max(0,this.activeProgress-1),this.ensureContainer()},ensureContainer(){document.getElementById("uno-notify-container")||a()}});const a=()=>{try{const t=document.createElement("div");t.id="uno-notify-container",t.setAttribute("x-data","{}");const i=`
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
                ${n()}
            </template>
        </div>

        <!-- Top Left -->
        <div class="fixed top-0 left-0 z-[100] p-4 flex flex-col space-y-2 items-start pointer-events-none">
            <template x-for="item in $store.notify.items.filter(i => i.position === 'top-left')" :key="item.id">
                ${n()}
            </template>
        </div>

        <!-- Bottom Right -->
        <div class="fixed bottom-0 right-0 z-[100] p-4 flex flex-col-reverse space-y-2 space-y-reverse items-end pointer-events-none">
            <template x-for="item in $store.notify.items.filter(i => i.position === 'bottom-right')" :key="item.id">
                ${n()}
            </template>
        </div>

        <!-- Bottom Left -->
        <div class="fixed bottom-0 left-0 z-[100] p-4 flex flex-col-reverse space-y-2 space-y-reverse items-start pointer-events-none">
            <template x-for="item in $store.notify.items.filter(i => i.position === 'bottom-left')" :key="item.id">
                ${n()}
            </template>
        </div>
      `;t.innerHTML=i,document.body.appendChild(t)}catch(t){x.error("Notify: Failed to inject container",t)}},n=()=>`
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
    `;document.body?a():document.addEventListener("DOMContentLoaded",a)};let k=0,L="";function _(){k===0&&(L=document.body.style.overflow||"",document.body.style.overflow="hidden"),k++}function D(){k=Math.max(0,k-1),k===0&&(document.body.style.overflow=L||"",L="")}async function N(s){const a=await fetch(s);return(a.headers.get("content-type")||"").includes("application/json")?a.json():a.text()}const W=s=>{s.store("box",{items:[],init(){s.effect(()=>{this.items.filter(t=>t.modal).length>0?_():D()})},async open(n,t={}){let i="dialog",e=t;typeof n=="object"?(e=n,i=e.type||"dialog"):i=n;const r=e.boxId||Math.random().toString(36).slice(2,10),l={id:r,type:i,title:e.title||"",content:e.content||"",width:e.width||400,height:e.height,modal:e.modal!==!1,overlayOpacity:e.overlayOpacity??.5,position:e.position||"center",showCloseButton:!0,okText:e.okButtonText,cancelText:e.cancelButtonText,onOk:e.onOk,onCancel:e.onCancel,showOk:i==="confirm"||i==="message"||i==="dialog"&&e.showOkButton,showCancel:i==="confirm"||i==="dialog"&&e.showCancelButton,loading:!1,visible:!0},d=s.store("i18n");if(l.okText||(l.okText=d?d.t("box.ok"):"OK"),l.cancelText||(l.cancelText=d?d.t("box.cancel"):"Cancel"),e.targetId){const o=document.getElementById(e.targetId);o?l.content=o.innerHTML:(x.warn("Box: Target element not found:",e.targetId),l.content='<div class="p-4 text-red-500">Error: Content target not found</div>')}else if(e.ajaxUrl){l.loading=!0,this.add(l);try{const o=await N(e.ajaxUrl),u=typeof o=="object"?`<pre class="bg-gray-100 p-3 rounded overflow-x-auto whitespace-pre-wrap">${JSON.stringify(o,null,2)}</pre>`:String(o),f=this.items.find(c=>c.id===r);f&&(f.content=u,f.loading=!1)}catch(o){x.error("Box: Ajax failed",o);const u=this.items.find(f=>f.id===r);u&&(u.content='<div class="text-red-500">Error loading content</div>',u.loading=!1)}return r}else e.iframeUrl&&(l.content=`<iframe frameborder="0" class="w-full h-full" src="${e.iframeUrl}"></iframe>`);return this.add(l),r},add(n){this.items.push(n),this.ensureContainer()},close(n){this.items=this.items.filter(t=>t.id!==n)},closeAll(){this.items=[]},getPositionClass(n){const t=n||"center",i="fixed inset-0 flex transition-opacity duration-300 p-4 pointer-events-none",e={center:"items-center justify-center","top-center":"items-start justify-center","top-left":"items-start justify-start","top-right":"items-start justify-end","bottom-center":"items-end justify-center","bottom-left":"items-end justify-start","bottom-right":"items-end justify-end","center-left":"items-center justify-start","center-right":"items-center justify-end"};return`${i} ${e[t]||e.center}`},ensureContainer(){document.getElementById("uno-box-container")||a()}});const a=()=>{try{const n=document.createElement("div");n.id="uno-box-container",n.setAttribute("x-data","{}");const t=`
        <template x-for="(item, index) in $store.box.items" :key="item.id">
            <div :id="'uno-box-outer-container' + item.id"
                 :class="$store.box.getPositionClass(item.position)"
                 :style="\`z-index: \${100 + index}\`"
                 role="dialog"
                 aria-modal="true">
                 
                 <!-- Overlay: Covers viewport, blocks interaction, does NOT close modal on click -->
                 <div :id="'uno-box-overlay' + item.id"
                      class="absolute inset-0 bg-black transition-opacity duration-300 pointer-events-auto"
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
                      class="relative bg-white dark:bg-gray-900 rounded-lg shadow-2xl ring-1 ring-black/5 overflow-hidden flex flex-col max-h-[90vh] transition-all duration-300 transform pointer-events-auto"
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
      `;n.innerHTML=t,document.body.appendChild(n)}catch(n){x.error("Box: Failed to inject container",n)}};document.body?a():document.addEventListener("DOMContentLoaded",a),document.addEventListener("click",n=>{const t=n.target.closest("[data-box-type]");if(!t)return;const i=t.getAttribute("data-box-type");if(!["dialog","confirm","message"].includes(i))return;n.preventDefault();const e=t.dataset,r={title:e.boxTitle,content:e.boxContent,width:e.boxWidth?parseInt(e.boxWidth):void 0,height:e.boxHeight?parseInt(e.boxHeight):void 0,ajaxUrl:e.boxAjaxUrl,iframeUrl:e.boxIframeUrl,targetId:e.boxTargetId,position:e.boxPosition,okButtonText:e.boxOkText,cancelButtonText:e.boxCancelText,overlayOpacity:e.boxOverlayOpacity?parseFloat(e.boxOverlayOpacity):void 0,showCloseButton:e.boxShowClose!=="false",boxId:e.boxId,cache:e.boxCache==="true"};e.boxModal==="false"&&(r.modal=!1),s.store("box").open(i,r)})},z={validname:/^[a-zA-Z0-9_-]+$/,en:/^[a-zA-Z]+$/,domain:new RegExp("^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\\.)+[A-Za-z]{2,6}$"),mobile:/^\+?[\d\s-]{7,}$/,int:/^-?\d+$/,email:/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/,date:/^(\d{4}-\d{2}-\d{2}|\d{2}\/\d{2}\/\d{4}|\d{2}\.\d{2}\.\d{4})$/,url:/^(https?:\/\/)?([\da-z.-]+)\.([a-z.-]{2,6})([/\w .-]*)*\/?$/i,pwd:/^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/,strongpwd:/^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/};function Y(s){const a=s.match(/^\/(.*?)\/([gimsuy]*)$/);return a?{pattern:a[1],flags:a[2]}:{pattern:s,flags:""}}function Z(s){try{return s?JSON.parse(s):[]}catch{return[]}}const V=s=>{const a={invalid:`url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 12 12' width='12' height='12' fill='none' stroke='%23ef4444'%3e%3ccircle cx='6' cy='6' r='4.5'/%3e%3cpath stroke-linecap='round' d='M5.8 3.6h.4L6 6.5zM6 8.2af.6.6 0 110-1.2.6.6 0 010 1.2z'/%3e%3c/svg%3e")`,valid:`url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 8'%3e%3cpath fill='%2322c55e' d='M2.3 6.73L.6 4.53c-.4-1.04.46-1.4 1.1-.8l1.1 1.4 3.4-3.8c.6-.63 1.6-.27 1.2.7l-4 4.6c-.43.5-.8.4-1.1.1z'/%3e%3c/svg%3e")`},n=(t,i)=>{const e=s.store("i18n");if(!e)return t;const r={required:"Required",pattern:"Invalid format",min:`At least ${i?.n} characters`,max:`Max ${i?.n} characters`,int:"Integer required",email:"Invalid email address",url:"Invalid URL"};return e.t(`valid.${t}`,i)||r[t]||t};s.store("valid",{callbacks:{},registerCallback(t,i){typeof i=="function"&&(this.callbacks[t]=i)}}),s.data("validForm",()=>({fields:{},init(){this.$el.noValidate=!0},register(t,i,e=[]){if(this.fields[t])return;const r=[...e];if((i.hasAttribute("required")||i.dataset.validRequired==="true")&&r.push({type:"required"}),(i.type==="email"||i.dataset.validType==="email")&&r.push("email"),(i.type==="url"||i.dataset.validType==="url")&&r.push("url"),i.dataset.validType){const l=i.dataset.validType;l!=="email"&&l!=="url"&&r.push(l)}i.dataset.validPattern&&r.push({type:"regexp",expression:i.dataset.validPattern}),this.fields[t]={el:i,rules:r,value:i.value,error:null,valid:!0,dirty:!1,touched:!1,validating:!1}},update(t,i){this.fields[t]&&(this.fields[t].value=i,this.fields[t].dirty=!0,this.validate(t))},touch(t){this.fields[t]&&(this.fields[t].touched=!0,this.validate(t))},async validate(t){const i=this.fields[t];if(!i)return!0;i.validating=!0,i.error=null;const e=(i.value||"").trim();let r=null;for(const l of i.rules){const d=l.type||l,o=typeof l=="object"?l:{};if(!e){if(d==="required"){r=o.msg||n("required");break}if(d==="emptyquit")break;continue}if(z[d])z[d].test(e)||(r=o.msg||n(d)||n("pattern"));else if(d==="regexp"&&o.expression){const{pattern:u,flags:f}=Y(o.expression);new RegExp(u,f||"i").test(e)||(r=o.msg||n("pattern"))}else if(d==="limit")o.min&&e.length<o.min&&(r=o.msg||n("min",{n:o.min})),o.max&&e.length>o.max&&(r=o.msg||n("max",{n:o.max}));else if(d==="range"){const u=parseFloat(e);o.min!==void 0&&u<o.min&&(r=o.msg||`Min ${o.min}`),o.max!==void 0&&u>o.max&&(r=o.msg||`Max ${o.max}`)}else if(d==="compareto"){const u=document.getElementById(o.destid);u&&e!==u.value&&(r=o.msg||"Not match")}else if(d==="tag"){const u=o.mark||",",f=e.split(u).filter(c=>c.trim());o.max&&f.length>o.max&&(r=o.msg||`Max ${o.max} items`)}if(r)break;if(d==="callback"&&o.fn){const u=s.store("valid").callbacks[o.fn];if(u)try{const f=await u(i.el.id,e);f===!1||typeof f=="object"&&f.valid===!1?r=(typeof f=="object"?f.message:"")||"Invalid":typeof f=="string"&&(r=f)}catch{}}if(r)break}return i.error=r,i.valid=!r,i.validating=!1,i.valid},async validateAll(){const t=Object.keys(this.fields);return(await Promise.all(t.map(e=>this.validate(e)))).every(e=>e)},async submit(t){if(t&&(t.preventDefault(),t.stopPropagation()),Object.keys(this.fields).forEach(e=>{this.fields[e].touched=!0}),await this.validateAll())this.$el.dispatchEvent(new CustomEvent("valid-submit",{bubbles:!0})),!this.$el.hasAttribute("hx-post")&&!this.$el.hasAttribute("hx-get")&&this.$el.submit();else{const e=Object.keys(this.fields).find(r=>!this.fields[r].valid);e&&this.fields[e].el&&this.fields[e].el.focus()}}})),s.directive("validate",(t,i,{effect:e,cleanup:r})=>{const l=t.getAttribute("name");if(!l)return;let d=[];const o=t.getAttribute("valid");o&&(d=Z(o)),e(()=>{const c=s.$data(t);c&&typeof c.register=="function"&&c.register(l,t,d)});const u=()=>{const c=s.$data(t);if(c&&c.update&&c.update(l,t.value),d.find(h=>(h.type==="limit"||h==="limit")&&h.show===!0)){const h=d.find(m=>m.type==="limit"||m==="limit");h&&h.max&&t.value.length>h.max&&(t.value=t.value.substring(0,h.max),c&&c.update&&c.update(l,t.value))}},f=()=>{const c=s.$data(t);c&&c.touch&&c.touch(l)};t.addEventListener("input",u),t.addEventListener("blur",f),e(()=>{const c=s.$data(t);if(c&&c.fields&&c.fields[l]){const h=c.fields[l],m=!!h.error,p=h.valid&&h.dirty&&!h.error;t.classList.toggle("is-invalid",m),t.classList.toggle("is-valid",p),t.classList.toggle("validating",h.validating),m?(t.classList.add("!border-red-500","pr-10","bg-no-repeat","bg-[right_0.75rem_center]","bg-[length:1rem_1rem]"),t.classList.remove("!border-green-500"),t.style.backgroundImage=a.invalid):p?(t.classList.add("!border-green-500","pr-10","bg-no-repeat","bg-[right_0.75rem_center]","bg-[length:1rem_1rem]"),t.classList.remove("!border-red-500"),t.style.backgroundImage=a.valid):(t.classList.remove("!border-red-500","!border-green-500","pr-10","bg-no-repeat","bg-[right_0.75rem_center]","bg-[length:1rem_1rem]"),t.style.backgroundImage="");const v=t.nextElementSibling;if(v&&v.tagName==="DIV"){const b=v.querySelector(".animate-spin")||t.parentNode.querySelector(".animate-spin");b&&(b.style.display=h.validating?"block":"")}}}),r(()=>{t.removeEventListener("input",u),t.removeEventListener("blur",f)})}),s.directive("valid-error",(t,{expression:i},{effect:e})=>{const r=i;t.style.display="none",e(()=>{const l=s.$data(t);if(l&&l.fields&&l.fields[r]){const o=l.fields[r].error;o?(t.textContent=o,t.style.display="block",t.classList.add("text-red-600","text-xs"),t.classList.remove("text-gray-400")):t.style.display="none"}})}),s.directive("valid-count",(t,{expression:i},{effect:e})=>{const r=i;e(()=>{const l=s.$data(t);if(l&&l.fields&&l.fields[r]){const d=l.fields[r],o=d.rules.find(u=>u.type==="limit"||u==="limit");if(o&&o.max){const u=(d.value||"").length,f=o.max-u,c=s.store("i18n");t.textContent=c&&c.t("charCount",{count:f})||`Remaining: ${f}`,t.style.display="block"}else t.style.display="none"}})})};const E=(s="",a={})=>({active:s,options:{hash:!0,manual:!1,...a},loaded:{},init(){this.options&&this.options.hash&&this.$nextTick(()=>{const n=window.location.hash.substring(1);if(n){const t=this.$el.querySelector(`[data-tab-hash="${n}"]`);if(t){const i=t.getAttribute("data-tab-target");i&&(this.active=i)}}}),this.options.hash&&window.addEventListener("hashchange",()=>{const n=window.location.hash.substring(1);if(n){const t=this.$el.querySelector(`[data-tab-hash="${n}"]`);if(t){const i=t.getAttribute("data-tab-target");i&&(this.active=i)}}})},isActive(n){return this.active===n},async open(n,t=null,i=null){if(this.active=n,this.options.hash&&i&&window.location.hash.substring(1)!==i&&history.pushState(null,null,"#"+i),t){const e=document.getElementById(n);if(e&&!this.loaded[n]){e.innerHTML=`<div class="animate-pulse p-4 space-y-2">
             <div class="h-4 bg-gray-200 rounded w-3/4"></div>
             <div class="h-4 bg-gray-200 rounded w-1/2"></div>
           </div>`;try{const r=await fetch(t);if(!r.ok)throw new Error("Network response was not ok");const l=await r.text();e.innerHTML=l,e.setAttribute("data-tab-loaded","true"),this.loaded[n]=!0,e.dispatchEvent(new CustomEvent("tab-loaded",{bubbles:!0}))}catch(r){e.innerHTML=`<div class="text-red-500 p-4">Error loading content: ${r.message}</div>`}}}},trigger(n,t=null,i=null){return{"@click.prevent":()=>this.open(n,t,i),":class":()=>this.isActive(n)?"current":"",role:"tab","aria-selected":()=>this.isActive(n).toString(),"data-tab-target":n,"data-tab-hash":i||null,"data-tab-url":t||null}},tablist:{role:"tablist","@keydown.right.prevent":"next()","@keydown.left.prevent":"prev()","@keydown.home.prevent":"first()","@keydown.end.prevent":"last()"},panel(n){return{"x-show":()=>this.isActive(n),role:"tabpanel",id:n}}});typeof window<"u"&&(window.unoTabs=E);const F=s=>{s.data("unoTabs",E)};window.Alpine=y,window.htmx=w,typeof window<"u"&&(window.unoTabs=E);const T=s=>{if(!s){x.error("Uno: Alpine is undefined, cannot register plugins");return}if(s.store&&s.store("tooltip")){x.info("Uno: Plugins already registered");return}x.info("Uno: Registering plugins"),S(s),P(s),U(s),W(s),V(s),F(s)};(()=>{const s=y||window.Alpine;s&&(x.info("Uno: Alpine found, registering"),T(s)),typeof window<"u"&&(window.deferLoadingAlpine=a=>{try{T(window.Alpine)}catch{}a()}),document.addEventListener("alpine:init",()=>{T(window.Alpine)})})()}));
