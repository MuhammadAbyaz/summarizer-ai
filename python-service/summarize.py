import os
import boto3
import google.generativeai as genai


def summarizer_handler(event, context):
    genai.configure(api_key=os.getenv("GEMINI_API_KEY"))
    model = genai.GenerativeModel("gemini-1.5-flash")
    s3_client = boto3.client("s3")
    boto3.set_stream_logger("botocore", level="DEBUG")
    object_name = event["Records"][0]["s3"]["object"]["key"]

    response = s3_client.get_object(Bucket=os.getenv("BUCKET_NAME"), Key=object_name)
    content = response["Body"].read()
    prompt = "Summarize this file"
    response = model.generate_content(
        [{"mime_type": "application/pdf", "data": content}, prompt]
    )

    return {"statusCode": 200, "body": f"{response.text}"}
