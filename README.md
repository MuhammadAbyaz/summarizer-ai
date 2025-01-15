# Summarizer-AI

**Summarizer-AI** is an API-based service that provides PDF summarization using serverless architecture. Built with Go and AWS CDK, it leverages modern cloud technologies to deliver efficient and scalable solutions.

## Features
- **`/upload`**: Endpoint to upload PDF files. Files are stored in an S3 bucket, triggering a Lambda function for processing.
- **`/summarize`**: Internal Lambda (not user-accessible) processes the uploaded PDF using the Gemini API and stores the summarized content in Supabase.
- **`/get-summary`**: Endpoint to fetch the summary using a unique file ID. Retrieves data from Supabase and returns it to the user.

## Architecture
- **API Gateway**: Handles all API requests.
- **3 Lambda Functions**:
  1. Handles file uploads to S3.
  2. Summarizes PDFs using the Gemini API and stores results.
  3. Fetches summaries from Supabase based on user requests.
- **S3 Trigger**: Activates the summarization Lambda on file upload.
- **Supabase**: Stores summarized results for quick retrieval.

## Architecture Diagram
<img src="https://7b98e0sbxv.ufs.sh/f/k9KWkJn7v0HsnfTfe8uKVLwikQlHdvEX95n0IDBqgS2rzteG"/>

## Highlights
This project demonstrates the power of:
- Event-driven architecture.
- Modern serverless technologies.
- Seamless integration of cloud services for efficient processing.

---
Feel free to contribute or report any issues!