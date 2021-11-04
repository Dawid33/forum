#!/bin/bash
(cd backend; sudo docker build -t backend/backend:latest .)
(cd frontend; sudo docker build -t forum-console/forum-console:latest .)