# create-thumbnails-lambda

This lambda function will be invoked once you put image to an Amazon S3 Bucket.

The function gets images from the bucket,
creates thumbnails,
and
put them to another bucket.

## environment variables

| name | default value | note |
| :---: | :---: | :--- |
| IMAGE_RATE | 0.5 | magnification rate |
| OBJECT_BUCKET_NAME_ORIGINAL | original.images.mububoki | Lambda will be invoked when you put images to this bucket.
| OBJECT_BUCKET_NAME_THUMBNAIL | thumbnail.images.mububoki | Lambda puts thumbnails to this bucket. |
