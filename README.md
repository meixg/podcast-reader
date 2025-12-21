# podcast-reader

A web service that let people download podcast and read it.

## usage

1. Open the webpage, past the xiaoyuzhou podcast url like: https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3
2. Click submit button, the server will:
    1. download the audio file of this podcast and
    2. process the audio file using Tecent Cloud service to get transcript.
    3. submit the transcript to LLM to get a briefing doc.
3. Wait and get the briefing doc.

Since the process time maybe very long, there will be a process id to let user retrive result.

## overall architecture

### front-end page

There is a simple front-end page using Vue + vite.

### Go server

go server do the audio processing part.

### db

so far we don't need a db, just store the processing info and result in a specified folder.

