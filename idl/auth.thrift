namespace go auth


struct ReqBody {
    1: string username
    2: string password
    3: string phone
    4: string email
    5: string grantType
    6: string scope
}

struct Msg {
	1: i8 code
	2: string msg
}

struct Response {
	1: Msg msg
	2: string data
}

service AuthService {
    Response auth(1: ReqBody req)(api.post="/v1/auth")
}