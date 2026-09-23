#!/bin/bash

if [ $1 == "habit" ]; then
  cd backend/services/habit-service
elif [ $1 == "social" ]; then
  cd backend/services/social-service
elif [ $1 == "stats" ]; then
  cd backend/services/stats-service
elif [ $1 == "user" ]; then
  cd backend/services/user-service
elif [ $1 == "auth" ]; then
  cd backend/auth-services-node
fi