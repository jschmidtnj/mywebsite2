# graphql api

currently using golang

## potential documentation

https://github.com/2fd/graphdoc#demos

## run in development

`go run *.go`

## add to production

app engine: `gcloud app deploy`

## queries

- `http://localhost:port/graphql?query={post(type:"blog",id:"id"){title content id author}}`
- `http://localhost:port/graphql?query=mutation{addPost(type:"blog",typetitle:"asdf",content:"asdf",author:"asdf"){title}}`
- `http://localhost:port/graphql?query=mutation{updatePost(type:"blog",id:"id",title:"test123",author:"asdf"){title views}}`
- `http://localhost:port/graphql?query=mutation{deletePost(type:"blog",id:"5cef23a99833f8037391e3c6"){title views}}`
- `http://localhost:port/graphql?query={posts(type:"blog",perpage:10,page:0,searchterm:"asdf",sort:"title",ascending:false){title content views id author date}}`
- `http://localhost:port/graphql?query={account{id email}}`
- `http://localhost:port/graphql?query={user(id:"id"){email id password}}`
- `http://localhost:port/graphql?query=mutation{deleteUser(id:"id"){id email}}`

## edit CORS

The file `cors.json` is needed to allow file downloads from Firebase Storage. To configure CORS, download the [gsutil utility](https://cloud.google.com/storage/docs/gsutil_install) on your computer (again, linux preferred), and run `gcloud init` to sign in. Then run `export BOTO_CONFIG=/dev/null` on linux to prevent a bug in the program. Finally, run `gsutil cors set rules/cors.json gs://<your-cloud-storage-bucket>` in the parent directory to add the CORS rules, obviously changing the command for your storage bucket.
