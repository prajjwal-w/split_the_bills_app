# split_the_bills_app
this is a backend for the split the bills app where you can split the bills among your friends, the backend is fully developed in golang, I will add the features time to time and improve it.

You can signup, login, create group, edit group, add expense in group, remove user from group, split expense with users

A. Authentication Routes
-------------------------
1. POST -- /users/signup -->Register user in the system
   body JSON:
       {
            "first_name":"",
            "last_name":"",
            "password":"",
            "phone":"",
            "email":"",
            "user_type":"USER" or "ADMIN"
       }

2. POST -- /users/login -->login
   body JSON:
      {
        "email":"",
        "password":""
      }

B. User Routes
---------------
1. GET -- /users --> Get all users (role based access only ADMIN)
2. GET -- /users/:user_id --> Get user By User_Id
3. DELETE -- /users/:user_id --> Delete user by using User_Id
4. PUT -- /users --> Reset Password
   body JSON:
      {
        "userId": ,
        "password": ""
      }


C. Group Routes
------------------
1. POST -- /groups --> Create groups
   body JSON:
    {
        "group_name": "",
        "description": "",
        "created_by" :        -- user_id
    }

2. POST -- /groups/addUser --> Add user in a group
   body JSON:
   {
        "group_id": 2003,
        "user_id": 10004
   }

3. GET -- /groups/:user_id --> Get all group in which user is present

4. DELETE -- /group/deleteUser --> Delete user from a group
   body JSON:
   {
    "group_id":,
    "user_id":
   }


D. Expense Routes
------------------
1. POST -- /expenses --> add expenses in a group
   body JSON:
   {
    "amount": ,
    "description":"",
    "paid_by_id": ,
    "splitwith":[
        {
            "user_id":,
            "amt_split":,
            "status" : "" --settled or unsettled
        },
        {
            "user_id":,
            "amt_split": ,
            "status" : ""
        }
    ]
   }

2. GET -- /expenses/:user_id -->get all unsettled expenses of user

3. PUT -- /expenses  -->settle up the expense
   body JSON:
   {
    "payer_id": ,
    "payee_id": ,
    "amt": ,
   }