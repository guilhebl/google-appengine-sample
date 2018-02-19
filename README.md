# GOOGLE app engine sample

A sample app which sets up a web server which receives requests and generates the outputs of a producer-consumer
model. The app builds each task generating 3 random numbers which then each worker will sum up. The consumer
will then aggregate and multiply all the generated sums from each producer rendering the final result to the HTML.

## How to use

1. Download and install App engine SDK [Google App Engine SDK for Go](https://cloud.google.com/appengine/docs/standard/go/download)

2. Run in DEV mode (localhost) - `goapp serve`

3. Deploy to app engine - `gcloud app deploy`


Output sample:

worker 1 - numbers: 98, 91, 46

worker 2 - numbers: 81, 81, 94

worker 3 - numbers: 51, 93, 13

worker 4 - numbers: 63, 61, 56

Total: 1700121600