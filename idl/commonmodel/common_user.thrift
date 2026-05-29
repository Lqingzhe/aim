namespace go kitexcommonmodel
//UserService
struct UserInfo{
    1:i64 UserID
    2:string UserName
    3:string Introduction
    4:i64 BirthdayYear
    5:i64 BirthdayMonth
    6:i64 BirthdayDay
}
struct UserLoginInfo{
    1:i64 UserID
    2:string Password
    3:string Salt
}
struct RemarkInfo{
    1:i64 UserID
    2:i64 GoalUserID
    3:string NickName
}