package dto

// struct untuk pagination untuk response data (banyak dan sedikit)
type Paging struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalRows  int `json:"totalRows"`
	TotalPages int `json:"totalPages"`
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// response datanya satu
type SingleResponse struct {
	Status Status `json:"status"`
	Data   any    `json:"data"`
}

type PagingResponse struct {
	Status Status `json:"status"`
	Data   []any  `json:"data"`
	Paging Paging `json:"paging"`
}

//response data banyak

// gambaran bentuk json
// {
// 	"status":{
// 		"code" : 200,
// 		"message" : "success create data"
// 	}
// 	"data" : {
//   "id": "77902bf3-1bb0-4462-b05b-88878afe4411",
//   "billDate": "2024-06-27T00:00:00Z",
//   "customer": {
//     "id": "aa9e3f3e-0750-4feb-9393-cce26c12ee2a",
//     "name": "Rizqi",
//     "phoneNumber": "089898989",
//     "address": "jakarta",
//     "createdAt": "2024-06-26T16:39:04.707238Z",
//     "updatedAt": "0001-01-01T00:00:00Z"
//   },
//   "user": {
//     "id": "e3c74272-1b06-4f62-8ed1-996787e28734",
//     "name": "Aul",
//     "email": "aul@email.com",
//     "username": "aul",
//     "password": "123",
//     "role": "admin",
//     "createdAt": "2024-06-26T16:38:29.094532Z",
//     "updatedAt": "0001-01-01T00:00:00Z"
//   },
//   "billDetails": [
//     {
//       "id": "e08689ef-ddfe-4089-aa63-0a20c6c816cf",
//       "billId": "",
//       "product": {
//         "id": "0f872706-cdd1-48ce-9289-014bdbc346e7",
//         "name": "laundry express",
//         "price": 20000,
//         "type": "kg",
//         "createdAt": "2024-06-26T16:38:48.945811Z",
//         "updatedAt": "0001-01-01T00:00:00Z"
//       },
//       "qty": 10,
//       "price": 20000,
//       "createdAt": "0001-01-01T00:00:00Z",
//       "updatedAt": "0001-01-01T00:00:00Z"
//     }
//   ],
//   "createdAt": "0001-01-01T00:00:00Z",
//   "updatedAt": "0001-01-01T00:00:00Z"
// }
// }
