FROM pytorch/pytorch:2.0.1-cuda11.7-cudnn8-runtime

# please download the model from https://huggingface.co and put it in the same directory as Dockerfile
COPY ./bge-reranker-base-onnx-o4 ./bge-reranker-base-onnx-o4
COPY ./bge-m3 ./bge-m3

COPY requirements.txt ./
COPY uvicorn_config.json ./
RUN python3 -m pip install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

COPY api.py Dockerfile ./

ENTRYPOINT python3 api.py
