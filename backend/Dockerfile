FROM python:3-alpine

RUN python -m pip install --upgrade pip
RUN pip3 install gTTS

ADD app /app
ENTRYPOINT [ "/app" ]
