namespace go hello


struct ReqBody {
    1: string name
    2: i32 type
    3: string email
}

struct Request {
	1: string data
	2: string message
	3: ReqBody reqBody
}

struct Msg {
	1: i32 code
	2: string msg
}

struct Response {
	1: Msg msg
	2: string data
}

service HelloService {
    Response echo(1: Request req)
    Response Get()
    Response GetByParams(1: string data, 2: string msg)(api.get="/api/v1/getByParams/:data")
    Response GetH(1: i32 id)(api.get="/api/v1/getH/:id")
    Response Post(1: Request req)(api.post="/api/v1/create")
}