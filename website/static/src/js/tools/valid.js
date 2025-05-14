/**
 * jQuery表单验证插件
 * 
 * 这个插件提供了一种简单而强大的方式来验证表单输入。
 * 它支持多种验证类型，包括内置验证和自定义回调验证。
 * 
 * @author nobody
 * @version 1.2.0
 * @license MIT
 * 
 * 使用说明：
 * 
 * 1. 在HTML中为需要验证的输入元素添加valid属性：
 *    <input type="text" name="username" valid='["required","validname"]'>
 * 
 * 2. 支持的验证类型：
 *    - required: 必填字段
 *    - validname: 有效的用户名（字母、数字、下划线、连字符）
 *    - en: 仅英文字母
 *    - domain: 有效的域名
 *    - mobile: 手机号码
 *    - int: 整数
 *    - email: 电子邮件地址
 *    - date: 日期
 *    - url: URL地址
 *    - pwd: 密码（至少8个字符，包含字母和数字）
 *    - strongpwd: 强密码（至少8个字符，包含字母、数字和特殊字符）
 *    - limit: 字符长度限制
 *    - range: 数值范围
 *    - compareto: 与其他字段比较
 *    - callback: 自定义回调函数
 *    - tag: 标签数量限制
 * 
 * 3. 自定义回调函数示例：
 *    $.fn.valid.registerCallback('checkUsername', function(id, value) {
 *      return value.length >= 3 && value.length <= 20;
 *    });
 * 
 *    然后在HTML中使用：
 *    <input type="text" name="username" valid='[{"type":"callback", "fn":"checkUsername", "msg":"用户名长度必须在3到20个字符之间"}]'>
 * 
 * 4. 手动触发验证：
 *    $('#myForm').valid('checkError');
 * 
 * 5. 显示字符计数：
 *    <input type="text" name="description" valid='[{"type":"limit", "max": 100, "show": true}]'>
 *
 * 6. 多语言支持：
 *    通过设置环境变量 IDIARY_SUPPORT_LANGS 来指定支持的语言
 *    $.fn.valid.setLanguage('en'); // 手动设置为英语
 *
 * 7. 页面语言自动检测：
 *    插件会自动从页面中寻找名为 IDIARY_CURRENT_LANG 的全局变量来确定当前语言
 *    如果该变量不存在，则使用浏览器的默认语言设置
 *    例如: <script>var IDIARY_CURRENT_LANG = 'en';</script>
 *
 * 8. 支持的语言列表：
 *    $.fn.valid.setLanguage() 支持的语言代码如下:
 *    - ar: Arabic (العربية)
 *    - de: German (Deutsch)
 *    - en: English
 *    - es: Spanish (Español)
 *    - fr: French (Français)
 *    - ja: Japanese (日本語)
 *    - ko: Korean (한국어)
 *    - ru: Russian (Русский)
 *    - zh-cn: Chinese Simplified (简体中文)
 *    - zh-tw: Chinese Traditional (繁體中文)
 * 
 * 注意：确保在使用此插件之前已经加载了jQuery库。
 */

(function ($) {
  // 私有函数：计算字符串长度（使用标准的字符串长度属性）
  const strLen = (str) => str.length;

  // 私有函数：按字符数截取字符串（使用标准的substring方法）
  const substrByChar = (str, charLength) => str.substring(0, charLength);

  // 回调函数注册对象
  const callbackRegistry = {};

  // 全局语言变量名
  const LANG_VAR_NAME = 'IDIARY_CURRENT_LANG';

  // 支持的语言列表
  const supportedLanguages = [
    'ar',
    'de',
    'en',
    'es',
    'fr',
    'ja',
    'ko',
    'ru',
    'zh-cn',
    'zh-tw',
  ];

  // 语言代码标准化（转为小写）
  const normalizeLanguage = (lang) => {
    return (lang || '').toLowerCase();
  };

  // 获取当前语言
  const detectLanguage = () => {
    // 检查页面全局变量
    if (typeof window[LANG_VAR_NAME] !== 'undefined') {
      return normalizeLanguage(window[LANG_VAR_NAME]);
    }

    // 检查浏览器语言设置
    const browserLang = navigator.language || navigator.userLanguage;
    if (browserLang) {
      const normalizedBrowserLang = normalizeLanguage(browserLang);

      // 检查完整的语言代码
      if (supportedLanguages.includes(normalizedBrowserLang)) {
        return normalizedBrowserLang;
      }

      // 特殊处理中文的情况
      if (normalizedBrowserLang.startsWith('zh')) {
        // 对于中文，如果包含-tw或-hk或-mo，则使用繁体中文，否则使用简体中文
        if (
          normalizedBrowserLang.includes('-tw') ||
          normalizedBrowserLang.includes('-hk') ||
          normalizedBrowserLang.includes('-mo')
        ) {
          return 'zh-tw';
        } else {
          return 'zh-cn';
        }
      }

      // 如果浏览器语言包含国家代码（如en-US），只取主语言部分
      const mainLang = normalizedBrowserLang.split('-')[0];
      if (supportedLanguages.includes(mainLang)) {
        return mainLang;
      }
    }

    // 默认返回英文
    return 'en';
  };

  // 当前语言，自动检测
  let currentLanguage = detectLanguage();

  // 多语言信息对象
  const i18n = {
    'zh-cn': {
      validname: '请输入有效的用户名（包含字母、数字、下划线和连字符）',
      en: '请只输入英文字母（A-Z, a-z）',
      domain: '请输入有效的域名（如：example.com）',
      mobile: '请输入有效的手机号码（如：+86 123 4567 8901 或 12345678901）',
      int: '请输入整数（如：42 或 -17）',
      email: '请输入有效的电子邮件地址（如：user@example.com）',
      date: '请输入有效的日期（格式：YYYY-MM-DD、DD/MM/YYYY 或 DD.MM.YYYY）',
      url: '请输入有效的URL地址（如：http://www.example.com 或 www.example.com）',
      pwd: '请输入有效的密码（至少8个字符，包含至少一个字母和一个数字）',
      strongpwd:
        '请输入有效的强密码（至少8个字符，包含至少一个字母、一个数字和一个特殊字符 @$!%*#?&）',
      required: '必须填写',
      regexp: '输入数据不合法',
      range: '输入数字应该在{min}到{max}之间',
      rangeInt: '请输入数字',
      compareto: '请输入正确的内容',
      callback: '输入内容无效',
      tag: '标签个数不能超过{max}个',
      limitMin: '长度必须大于等于{min}个字符',
      limitMax: '长度必须小于等于{max}个字符',
      charCount: '可输入{count}个字符',
    },
    'zh-tw': {
      validname: '請輸入有效的用戶名（包含字母、數字、底線和連字符）',
      en: '請只輸入英文字母（A-Z, a-z）',
      domain: '請輸入有效的域名（如：example.com）',
      mobile: '請輸入有效的手機號碼（如：+886 123 4567 8901 或 12345678901）',
      int: '請輸入整數（如：42 或 -17）',
      email: '請輸入有效的電子郵件地址（如：user@example.com）',
      date: '請輸入有效的日期（格式：YYYY-MM-DD、DD/MM/YYYY 或 DD.MM.YYYY）',
      url: '請輸入有效的URL地址（如：http://www.example.com 或 www.example.com）',
      pwd: '請輸入有效的密碼（至少8個字符，包含至少一個字母和一個數字）',
      strongpwd:
        '請輸入有效的強密碼（至少8個字符，包含至少一個字母、一個數字和一個特殊字符 @$!%*#?&）',
      required: '必須填寫',
      regexp: '輸入數據不合法',
      range: '輸入數字應該在{min}到{max}之間',
      rangeInt: '請輸入數字',
      compareto: '請輸入正確的內容',
      callback: '輸入內容無效',
      tag: '標籤個數不能超過{max}個',
      limitMin: '長度必須大於等於{min}個字符',
      limitMax: '長度必須小於等於{max}個字符',
      charCount: '可輸入{count}個字符',
    },
    en: {
      validname:
        'Please enter a valid username (letters, numbers, underscores, and hyphens)',
      en: 'Please enter only English letters (A-Z, a-z)',
      domain: 'Please enter a valid domain name (e.g.: example.com)',
      mobile:
        'Please enter a valid phone number (e.g.: +86 123 4567 8901 or 12345678901)',
      int: 'Please enter an integer (e.g.: 42 or -17)',
      email: 'Please enter a valid email address (e.g.: user@example.com)',
      date: 'Please enter a valid date (format: YYYY-MM-DD, DD/MM/YYYY, or DD.MM.YYYY)',
      url: 'Please enter a valid URL (e.g.: http://www.example.com or www.example.com)',
      pwd: 'Please enter a valid password (at least 8 characters, including at least one letter and one number)',
      strongpwd:
        'Please enter a valid strong password (at least 8 characters, including at least one letter, one number, and one special character @$!%*#?&)',
      required: 'This field is required',
      regexp: 'Input data is invalid',
      range: 'The number should be between {min} and {max}',
      rangeInt: 'Please enter a number',
      compareto: 'Please enter the correct content',
      callback: 'Input content is invalid',
      tag: 'Number of tags cannot exceed {max}',
      limitMin: 'Length must be at least {min} characters',
      limitMax: 'Length must not exceed {max} characters',
      charCount: '{count} characters left',
    },
    ar: {
      validname:
        'الرجاء إدخال اسم مستخدم صالح (أحرف وأرقام وشرطات سفلية وواصلات)',
      en: 'الرجاء إدخال أحرف إنجليزية فقط (A-Z, a-z)',
      domain: 'الرجاء إدخال اسم نطاق صالح (مثال: example.com)',
      mobile:
        'الرجاء إدخال رقم هاتف صالح (مثال: ‎+962 123 4567 890 أو 1234567890)',
      int: 'الرجاء إدخال عدد صحيح (مثال: 42 أو -17)',
      email: 'الرجاء إدخال عنوان بريد إلكتروني صالح (مثال: user@example.com)',
      date: 'الرجاء إدخال تاريخ صالح (بالتنسيق: YYYY-MM-DD أو DD/MM/YYYY أو DD.MM.YYYY)',
      url: 'الرجاء إدخال رابط صالح (مثال: http://www.example.com أو www.example.com)',
      pwd: 'الرجاء إدخال كلمة مرور صالحة (8 أحرف على الأقل، تتضمن حرفاً واحداً ورقماً واحداً على الأقل)',
      strongpwd:
        'الرجاء إدخال كلمة مرور قوية (8 أحرف على الأقل، تتضمن حرفاً واحداً ورقماً واحداً وحرفاً خاصاً واحداً @$!%*#?& على الأقل)',
      required: 'هذا الحقل مطلوب',
      regexp: 'البيانات المدخلة غير صالحة',
      range: 'يجب أن يكون الرقم بين {min} و {max}',
      rangeInt: 'الرجاء إدخال رقم',
      compareto: 'الرجاء إدخال المحتوى الصحيح',
      callback: 'المحتوى المدخل غير صالح',
      tag: 'عدد العلامات لا يمكن أن يتجاوز {max}',
      limitMin: 'يجب أن يكون الطول {min} أحرف على الأقل',
      limitMax: 'يجب ألا يتجاوز الطول {max} حرفاً',
      charCount: 'متبقي {count} حرف',
    },
    de: {
      validname:
        'Bitte geben Sie einen gültigen Benutzernamen ein (Buchstaben, Zahlen, Unterstriche und Bindestriche)',
      en: 'Bitte geben Sie nur englische Buchstaben ein (A-Z, a-z)',
      domain:
        'Bitte geben Sie einen gültigen Domainnamen ein (z.B.: example.com)',
      mobile:
        'Bitte geben Sie eine gültige Telefonnummer ein (z.B.: +49 123 4567890 oder 12345678901)',
      int: 'Bitte geben Sie eine ganze Zahl ein (z.B.: 42 oder -17)',
      email:
        'Bitte geben Sie eine gültige E-Mail-Adresse ein (z.B.: user@example.com)',
      date: 'Bitte geben Sie ein gültiges Datum ein (Format: YYYY-MM-DD, DD/MM/YYYY oder DD.MM.YYYY)',
      url: 'Bitte geben Sie eine gültige URL ein (z.B.: http://www.example.com oder www.example.com)',
      pwd: 'Bitte geben Sie ein gültiges Passwort ein (mindestens 8 Zeichen, einschließlich mindestens eines Buchstabens und einer Zahl)',
      strongpwd:
        'Bitte geben Sie ein starkes Passwort ein (mindestens 8 Zeichen, einschließlich mindestens eines Buchstabens, einer Zahl und eines Sonderzeichens @$!%*#?&)',
      required: 'Dieses Feld ist erforderlich',
      regexp: 'Eingabedaten sind ungültig',
      range: 'Die Zahl sollte zwischen {min} und {max} liegen',
      rangeInt: 'Bitte geben Sie eine Zahl ein',
      compareto: 'Bitte geben Sie den richtigen Inhalt ein',
      callback: 'Eingabeinhalt ist ungültig',
      tag: 'Die Anzahl der Tags darf {max} nicht überschreiten',
      limitMin: 'Die Länge muss mindestens {min} Zeichen betragen',
      limitMax: 'Die Länge darf {max} Zeichen nicht überschreiten',
      charCount: 'Noch {count} Zeichen übrig',
    },
    es: {
      validname:
        'Por favor, introduzca un nombre de usuario válido (letras, números, guiones bajos y guiones)',
      en: 'Por favor, introduzca solo letras inglesas (A-Z, a-z)',
      domain:
        'Por favor, introduzca un nombre de dominio válido (ej.: example.com)',
      mobile:
        'Por favor, introduzca un número de teléfono válido (ej.: +34 123 456 789 o 123456789)',
      int: 'Por favor, introduzca un número entero (ej.: 42 o -17)',
      email:
        'Por favor, introduzca una dirección de correo electrónico válida (ej.: usuario@example.com)',
      date: 'Por favor, introduzca una fecha válida (formato: YYYY-MM-DD, DD/MM/YYYY o DD.MM.YYYY)',
      url: 'Por favor, introduzca una URL válida (ej.: http://www.example.com o www.example.com)',
      pwd: 'Por favor, introduzca una contraseña válida (al menos 8 caracteres, incluyendo al menos una letra y un número)',
      strongpwd:
        'Por favor, introduzca una contraseña fuerte válida (al menos 8 caracteres, incluyendo al menos una letra, un número y un carácter especial @$!%*#?&)',
      required: 'Este campo es obligatorio',
      regexp: 'Los datos introducidos no son válidos',
      range: 'El número debe estar entre {min} y {max}',
      rangeInt: 'Por favor, introduzca un número',
      compareto: 'Por favor, introduzca el contenido correcto',
      callback: 'El contenido introducido no es válido',
      tag: 'El número de etiquetas no puede exceder {max}',
      limitMin: 'La longitud debe ser de al menos {min} caracteres',
      limitMax: 'La longitud no debe exceder {max} caracteres',
      charCount: 'Quedan {count} caracteres',
    },
    fr: {
      validname:
        "Veuillez entrer un nom d'utilisateur valide (lettres, chiffres, tirets bas et tirets)",
      en: 'Veuillez entrer uniquement des lettres anglaises (A-Z, a-z)',
      domain: 'Veuillez entrer un nom de domaine valide (ex : example.com)',
      mobile:
        'Veuillez entrer un numéro de téléphone valide (ex : +33 1 23 45 67 89 ou 0123456789)',
      int: 'Veuillez entrer un nombre entier (ex : 42 ou -17)',
      email:
        'Veuillez entrer une adresse e-mail valide (ex : utilisateur@example.com)',
      date: 'Veuillez entrer une date valide (format : YYYY-MM-DD, DD/MM/YYYY ou DD.MM.YYYY)',
      url: 'Veuillez entrer une URL valide (ex : http://www.example.com ou www.example.com)',
      pwd: 'Veuillez entrer un mot de passe valide (au moins 8 caractères, incluant au moins une lettre et un chiffre)',
      strongpwd:
        'Veuillez entrer un mot de passe fort valide (au moins 8 caractères, incluant au moins une lettre, un chiffre et un caractère spécial @$!%*#?&)',
      required: 'Ce champ est obligatoire',
      regexp: 'Les données saisies sont invalides',
      range: 'Le nombre doit être compris entre {min} et {max}',
      rangeInt: 'Veuillez entrer un nombre',
      compareto: 'Veuillez entrer le contenu correct',
      callback: 'Le contenu saisi est invalide',
      tag: 'Le nombre de tags ne peut pas dépasser {max}',
      limitMin: "La longueur doit être d'au moins {min} caractères",
      limitMax: 'La longueur ne doit pas dépasser {max} caractères',
      charCount: '{count} caractères restants',
    },
    ja: {
      validname:
        '有効なユーザー名を入力してください（文字、数字、アンダースコア、ハイフン）',
      en: '英字のみを入力してください (A-Z, a-z)',
      domain: '有効なドメイン名を入力してください（例：example.com）',
      mobile:
        '有効な電話番号を入力してください（例：+81 90 1234 5678 または 09012345678）',
      int: '整数を入力してください（例：42 または -17）',
      email: '有効なメールアドレスを入力してください（例：user@example.com）',
      date: '有効な日付を入力してください（形式：YYYY-MM-DD、DD/MM/YYYY、または DD.MM.YYYY）',
      url: '有効なURLを入力してください（例：http://www.example.com または www.example.com）',
      pwd: '有効なパスワードを入力してください（少なくとも8文字、少なくとも1つの文字と1つの数字を含む）',
      strongpwd:
        '有効な強力なパスワードを入力してください（少なくとも8文字、少なくとも1つの文字、1つの数字、1つの特殊文字 @$!%*#?& を含む）',
      required: 'このフィールドは必須です',
      regexp: '入力データが無効です',
      range: '数値は {min} から {max} の間である必要があります',
      rangeInt: '数値を入力してください',
      compareto: '正しい内容を入力してください',
      callback: '入力内容が無効です',
      tag: 'タグの数は {max} を超えることはできません',
      limitMin: '長さは少なくとも {min} 文字である必要があります',
      limitMax: '長さは {max} 文字を超えてはなりません',
      charCount: 'あと {count} 文字入力できます',
    },
    ko: {
      validname: '유효한 사용자 이름을 입력하세요 (문자, 숫자, 밑줄 및 하이픈)',
      en: '영문 알파벳만 입력하세요 (A-Z, a-z)',
      domain: '유효한 도메인 이름을 입력하세요 (예: example.com)',
      mobile:
        '유효한 전화번호를 입력하세요 (예: +82 10 1234 5678 또는 01012345678)',
      int: '정수를 입력하세요 (예: 42 또는 -17)',
      email: '유효한 이메일 주소를 입력하세요 (예: user@example.com)',
      date: '유효한 날짜를 입력하세요 (형식: YYYY-MM-DD, DD/MM/YYYY 또는 DD.MM.YYYY)',
      url: '유효한 URL을 입력하세요 (예: http://www.example.com 또는 www.example.com)',
      pwd: '유효한 비밀번호를 입력하세요 (최소 8자, 최소 하나의 문자와 하나의 숫자 포함)',
      strongpwd:
        '유효한 강력한 비밀번호를 입력하세요 (최소 8자, 최소 하나의 문자, 하나의 숫자 및 하나의 특수 문자 @$!%*#?& 포함)',
      required: '이 필드는 필수입니다',
      regexp: '입력 데이터가 유효하지 않습니다',
      range: '숫자는 {min}에서 {max} 사이여야 합니다',
      rangeInt: '숫자를 입력하세요',
      compareto: '올바른 내용을 입력하세요',
      callback: '입력 내용이 유효하지 않습니다',
      tag: '태그 수는 {max}개를 초과할 수 없습니다',
      limitMin: '길이는 최소 {min}자 이상이어야 합니다',
      limitMax: '길이는 {max}자를 초과할 수 없습니다',
      charCount: '{count}자 남았습니다',
    },
    ru: {
      validname:
        'Пожалуйста, введите действительное имя пользователя (буквы, цифры, подчеркивания и дефисы)',
      en: 'Пожалуйста, введите только английские буквы (A-Z, a-z)',
      domain:
        'Пожалуйста, введите действительное доменное имя (например: example.com)',
      mobile:
        'Пожалуйста, введите действительный номер телефона (например: +7 123 456 7890 или 12345678901)',
      int: 'Пожалуйста, введите целое число (например: 42 или -17)',
      email:
        'Пожалуйста, введите действительный адрес электронной почты (например: user@example.com)',
      date: 'Пожалуйста, введите действительную дату (формат: YYYY-MM-DD, DD/MM/YYYY или DD.MM.YYYY)',
      url: 'Пожалуйста, введите действительный URL (например: http://www.example.com или www.example.com)',
      pwd: 'Пожалуйста, введите действительный пароль (минимум 8 символов, включая как минимум одну букву и одну цифру)',
      strongpwd:
        'Пожалуйста, введите действительный надежный пароль (минимум 8 символов, включая как минимум одну букву, одну цифру и один специальный символ @$!%*#?&)',
      required: 'Это поле обязательно для заполнения',
      regexp: 'Введенные данные недействительны',
      range: 'Число должно быть между {min} и {max}',
      rangeInt: 'Пожалуйста, введите число',
      compareto: 'Пожалуйста, введите правильное содержание',
      callback: 'Введенное содержание недействительно',
      tag: 'Количество тегов не может превышать {max}',
      limitMin: 'Длина должна быть не менее {min} символов',
      limitMax: 'Длина не должна превышать {max} символов',
      charCount: 'Осталось {count} символов',
    },
  };

  // 获取消息
  const getMessage = (key, params = {}) => {
    const messages = i18n[currentLanguage] || i18n['en'];
    let message = messages[key] || i18n['en'][key] || key;

    // 替换参数
    Object.keys(params).forEach((param) => {
      message = message.replace(`{${param}}`, params[param]);
    });

    return message;
  };

  $.fn.valid = function () {
    const validators = {
      validname: {
        regex: /^[a-zA-Z0-9_-]+$/,
        key: 'validname',
      },
      en: {
        regex: /^[a-zA-Z]+$/,
        key: 'en',
      },
      domain: {
        regex: /^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\.)+[A-Za-z]{2,6}$/,
        key: 'domain',
      },
      mobile: {
        regex: /^\+?[\d\s-]{7,}$/,
        key: 'mobile',
      },
      int: {
        regex: /^-?\d+$/,
        key: 'int',
      },
      email: {
        regex:
          /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/,
        key: 'email',
      },
      date: {
        regex: /^(\d{4}-\d{2}-\d{2}|\d{2}\/\d{2}\/\d{4}|\d{2}\.\d{2}\.\d{4})$/,
        key: 'date',
      },
      url: {
        regex:
          /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/i,
        key: 'url',
      },
      pwd: {
        regex: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/,
        key: 'pwd',
      },
      strongpwd: {
        regex: /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/,
        key: 'strongpwd',
      },
    };

    const methods = {
      /**
       * 初始化验证
       */
      init() {
        this.find('[valid]').each((_, obj) => {
          const $obj = $(obj);
          if ($obj.data('loadedvalid')) return;

          $obj.after('<div />').data('loadedvalid', true);
          $obj.on('input blur', () => methods.checkElement($obj));

          const validRules = JSON.parse($obj.attr('valid'));
          validRules.forEach((json) => {
            if (json.show === true) {
              $obj.after('<span />');
              $obj.on('input', methods.updateCharCount);
              methods.updateCharCount.call($obj[0]);
            }
          });
        });
      },

      /**
       * 更新字符计数
       */
      updateCharCount() {
        const $obj = $(this);
        const value = $obj.val();
        const validRules = JSON.parse($obj.attr('valid'));

        for (const json of validRules) {
          if (json.show === true) {
            const maxlen = json.max;
            if (strLen(value) > maxlen) {
              $obj.val(substrByChar(value, maxlen));
            }
            const remainingChars = maxlen - strLen($obj.val());
            $obj
              .nextAll('span:first')
              .html(getMessage('charCount', { count: remainingChars }));
            break;
          }
        }
      },

      /**
       * 检查所有元素的错误
       * @returns {boolean} 是否有错误
       */
      checkError() {
        const errors = Array.from($(this).find('[valid]')).map((obj) =>
          methods.checkElement($(obj)),
        );
        return errors.some(Boolean);
      },

      /**
       * 检查单个元素的验证
       * @param {jQuery} $obj - 要检查的jQuery对象
       * @returns {boolean} 是否有错误
       */
      checkElement($obj) {
        const validRules = JSON.parse($obj.attr('valid'));
        const $labelError = $obj.nextAll('div:first');
        let hasError = false;
        let errorMsg = '';

        for (const json of validRules) {
          const type = json.type || json;
          const val = $obj.val();

          switch (type) {
            case 'emptyquit':
              if (strLen(val) === 0) return false;
              break;
            case 'required':
              if (strLen(val) === 0)
                errorMsg = json.msg || getMessage('required');
              break;
            case 'limit':
              const len = strLen(val);
              if (json.min !== undefined && len < json.min) {
                errorMsg =
                  json.msg || getMessage('limitMin', { min: json.min });
              }
              if (json.max !== undefined && len > json.max) {
                errorMsg =
                  json.msg || getMessage('limitMax', { max: json.max });
              }
              break;
            case 'validname':
            case 'en':
            case 'domain':
            case 'mobile':
            case 'int':
            case 'email':
            case 'date':
            case 'url':
              if (
                validators[type] &&
                !validators[type].regex.test(val) &&
                val !== ''
              ) {
                errorMsg = json.msg || getMessage(validators[type].key);
              }
              break;
            case 'regexp':
              const reg = new RegExp(`^${json.expression}$`, 'iu');
              if (!reg.test(val)) errorMsg = json.msg || getMessage('regexp');
              break;
            case 'range':
              if (/^\d+$/.test(val)) {
                const numVal = parseInt(val, 10);
                if (numVal < json.min || numVal > json.max) {
                  errorMsg =
                    json.msg ||
                    getMessage('range', { min: json.min, max: json.max });
                }
              } else {
                errorMsg = getMessage('rangeInt');
              }
              break;
            case 'compareto':
              if ($(`#${json.destid}`).val() !== val || val === '') {
                errorMsg = json.msg || getMessage('compareto');
              }
              break;
            case 'callback':
              if (typeof callbackRegistry[json.fn] === 'function') {
                if (callbackRegistry[json.fn]($obj.attr('id'), val) === false) {
                  errorMsg = json.msg || getMessage('callback');
                }
              } else {
                console.error(`Callback function ${json.fn} is not registered`);
              }
              break;
            case 'tag':
              if (val.split(json.mark || ' ').length > json.max && val !== '') {
                errorMsg = json.msg || getMessage('tag', { max: json.max });
              }
              break;
          }

          if (errorMsg) {
            $labelError.text(errorMsg).attr('class', 'invalid-feedback');
            $obj.removeClass('is-invalid is-valid').addClass('is-invalid');
            hasError = true;
            break;
          }
        }

        if (!hasError) {
          $labelError.text('');
          $obj.removeClass('is-invalid is-valid').addClass('is-valid');
        }
        return hasError;
      },
    };

    if (typeof arguments[0] === 'string' && methods[arguments[0]]) {
      return methods[arguments[0]].apply(
        this,
        Array.prototype.slice.call(arguments, 1),
      );
    } else if (typeof arguments[0] === 'object' || !arguments[0]) {
      return methods.init.apply(this, arguments);
    } else {
      $.error('Method ' + arguments[0] + ' does not exist on jQuery.valid');
    }
  };

  /**
   * 注册自定义回调函数
   * @param {string} name - 回调函数的名称
   * @param {function} callback - 回调函数
   */
  $.fn.valid.registerCallback = function (name, callback) {
    if (typeof callback === 'function') {
      callbackRegistry[name] = callback;
    } else {
      console.error(`Invalid callback function for ${name}`);
    }
  };

  /**
   * 设置当前语言
   * @param {string} language - 语言代码
   */
  $.fn.valid.setLanguage = function (language) {
    const normalizedLang = normalizeLanguage(language);
    if (i18n[normalizedLang]) {
      currentLanguage = normalizedLang;
      return true;
    } else {
      console.warn(
        `Language ${language} is not supported, using default language`,
      );
      return false;
    }
  };

  /**
   * 添加新的语言支持
   * @param {string} language - 语言代码
   * @param {object} messages - 消息对象
   */
  $.fn.valid.addLanguage = function (language, messages) {
    const normalizedLang = normalizeLanguage(language);
    i18n[normalizedLang] = { ...i18n['en'], ...messages };

    // 如果新添加的语言不在支持列表中，则添加它
    if (!supportedLanguages.includes(normalizedLang)) {
      supportedLanguages.push(normalizedLang);
    }
  };

  /**
   * 获取当前语言
   * @returns {string} 当前语言代码
   */
  $.fn.valid.getLanguage = function () {
    return currentLanguage;
  };

  /**
   * 获取支持的语言列表
   * @returns {array} 支持的语言代码列表
   */
  $.fn.valid.getSupportedLanguages = function () {
    return supportedLanguages.slice();
  };

  // 使用 jQuery 的简写语法和方法链
  $(() => $('form').valid());
})(jQuery);
