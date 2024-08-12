package errors

// 定义所有的错误码
var (
    InputParamErr            = NewError(100000, "input_param_error")
    RegisterUserNameExists   = NewError(100001, "register_user_name_exists")
    RegisterUserEmailExists  = NewError(100002, "register_user_email_exists")
    ParseInputParamErr       = NewError(100003, "parse_input_param_error")
    UploadFileTypeErr        = NewError(100004, "upload_file_type_error")
    UploadFileTooLarge       = NewError(100005, "upload_file_too_large")
    UploadFileErr            = NewError(100006, "upload_file_error")
    UploadFileEmpty          = NewError(100007, "upload_file_empty")

    BizErr          = NewError(200000, "biz_error")
    UserNotExists   = NewError(200001, "user_not_exists")
    UserPwdErr      = NewError(200002, "user_pwd_error")
    ParseJsonErr    = NewError(200003, "parse_json_error")
    MaxDaybookLimitErr = NewError(200004, "max_daybook_limit_error")

    DbOperationErr  = NewError(300000, "db_operation_error")
    ExternalErr     = NewError(400000, "external_error")
    InternalErr     = NewError(500000, "internal_error")
)