# Comprehensive QA Testing Guide for AI Resume Builder

This document provides systematic QA testing procedures for all functionality in the AI Resume Builder backend.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Authentication Testing](#authentication-testing)
3. [Profile Management Testing](#profile-management-testing)
4. [Experience Management Testing](#experience-management-testing)
5. [Education Management Testing](#education-management-testing)
6. [Skills Management Testing](#skills-management-testing)
7. [Document Generation Testing](#document-generation-testing)
8. [Document Management Testing](#document-management-testing)
9. [Authorization & Security Testing](#authorization--security-testing)
10. [Error Handling Testing](#error-handling-testing)
11. [Database Integrity Testing](#database-integrity-testing)
12. [Performance Testing](#performance-testing)
13. [Bug Tracking Template](#bug-tracking-template)

---

## Prerequisites

### 1. Environment Setup

```bash
# Start the database
make docker-up

# Apply migrations
make migrate-up

# Start the server
make run
```

### 2. Health Check

```bash
# Verify server is running
curl http://localhost:8080/api/v1/health

# Expected response:
# {
#   "status": "ok",
#   "database": "connected"
# }
```

---

## Authentication Testing

### Test Case 1.1: User Registration - Valid Data

**Objective:** Verify successful user registration with valid credentials

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "qa.test1@example.com",
    "password": "SecurePass123!",
    "full_name": "QA Test User"
  }'
```

**Expected Result:**
- Status: 201 Created
- Response includes `access_token` and `refresh_token`
- User created with `free_generations_left: 2`
- Empty profile automatically created
- Password is hashed (never returned)

**Verification:**
```sql
SELECT email, full_name, free_generations_left, is_premium, password_hash IS NOT NULL as has_password
FROM users WHERE email = 'qa.test1@example.com';

SELECT user_id FROM profiles WHERE user_id = (SELECT id FROM users WHERE email = 'qa.test1@example.com');
```

### Test Case 1.2: User Registration - Duplicate Email

**Objective:** Verify duplicate email rejection

```bash
# Register same email twice
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "qa.test1@example.com",
    "password": "AnotherPass456!",
    "full_name": "Another User"
  }'
```

**Expected Result:**
- Status: 409 Conflict
- Error code: `EMAIL_EXISTS`
- No new user created

### Test Case 1.3: User Registration - Weak Password

**Objective:** Verify password strength validation

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "weakpass@example.com",
    "password": "weak",
    "full_name": "Weak Password User"
  }'
```

**Expected Result:**
- Status: 400 Bad Request
- Error code: `INVALID_REQUEST`
- Message mentions password requirements (minimum 8 characters)

### Test Case 1.4: User Registration - Invalid Email Format

**Objective:** Verify email validation

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "not-an-email",
    "password": "SecurePass123!",
    "full_name": "Invalid Email User"
  }'
```

**Expected Result:**
- Status: 400 Bad Request
- Error indicates invalid email format

### Test Case 1.5: User Login - Valid Credentials

**Objective:** Verify successful login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "qa.test1@example.com",
    "password": "SecurePass123!"
  }'
```

**Expected Result:**
- Status: 200 OK
- Response includes `access_token` and `refresh_token`
- User data returned (without password)

**Save Token for Next Tests:**
```bash
export QA_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"qa.test1@example.com","password":"SecurePass123!"}' \
  | jq -r '.access_token')

echo "Token: $QA_TOKEN"
```

### Test Case 1.6: User Login - Invalid Password

**Objective:** Verify incorrect password rejection

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "qa.test1@example.com",
    "password": "WrongPassword123!"
  }'
```

**Expected Result:**
- Status: 401 Unauthorized
- Error code: `INVALID_CREDENTIALS`

### Test Case 1.7: User Login - Non-existent User

**Objective:** Verify login with non-existent email

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "nonexistent@example.com",
    "password": "AnyPassword123!"
  }'
```

**Expected Result:**
- Status: 401 Unauthorized
- Error code: `INVALID_CREDENTIALS`

### Test Case 1.8: Token Refresh - Valid Token

**Objective:** Verify token refresh mechanism

```bash
# First, get refresh token
REFRESH_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"qa.test1@example.com","password":"SecurePass123!"}' \
  | jq -r '.refresh_token')

# Use refresh token to get new access token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}"
```

**Expected Result:**
- Status: 200 OK
- New `access_token` returned
- New `refresh_token` returned

### Test Case 1.9: Token Refresh - Invalid Token

**Objective:** Verify rejection of invalid refresh token

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"invalid.token.here"}'
```

**Expected Result:**
- Status: 401 Unauthorized
- Error code: `INVALID_TOKEN`

### Test Case 1.10: Logout

**Objective:** Verify logout functionality

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $QA_TOKEN"
```

**Expected Result:**
- Status: 200 OK
- Success message returned
- Note: Actual token invalidation is client-side

---

## Profile Management Testing

### Test Case 2.1: Get Profile - Empty Profile

**Objective:** Verify retrieval of auto-created empty profile

```bash
curl http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer $QA_TOKEN"
```

**Expected Result:**
- Status: 200 OK
- Profile exists with user's email and name
- Optional fields are null (phone, location, linkedin_url, etc.)

### Test Case 2.2: Update Profile - Valid Data

**Objective:** Verify successful profile update

```bash
curl -X PUT http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+1234567890",
    "location": "San Francisco, CA",
    "linkedin_url": "https://linkedin.com/in/qatest",
    "github_url": "https://github.com/qatest",
    "website_url": "https://qatest.dev",
    "summary": "Experienced QA engineer with expertise in automated testing"
  }'
```

**Expected Result:**
- Status: 200 OK
- All fields updated correctly
- `updated_at` timestamp changed

**Verification:**
```sql
SELECT phone, location, linkedin_url, summary, updated_at
FROM profiles
WHERE user_id = (SELECT id FROM users WHERE email = 'qa.test1@example.com');
```

### Test Case 2.3: Update Profile - Partial Update

**Objective:** Verify partial profile updates don't clear other fields

```bash
curl -X PUT http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "summary": "Updated summary only"
  }'
```

**Expected Result:**
- Status: 200 OK
- Only summary updated
- Other fields (phone, location, etc.) remain unchanged

### Test Case 2.4: Get Profile - Populated Profile

**Objective:** Verify retrieval of updated profile

```bash
curl http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer $QA_TOKEN"
```

**Expected Result:**
- Status: 200 OK
- All previously set fields returned correctly

---

## Experience Management Testing

### Test Case 3.1: Create Experience - Valid Data

**Objective:** Verify successful experience creation

```bash
curl -X POST http://localhost:8080/api/v1/profile/experience \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company": "Tech Corp",
    "position": "Senior QA Engineer",
    "start_date": "2020-01-15",
    "end_date": "2023-06-30",
    "is_current": false,
    "description": "Led QA automation team for enterprise applications",
    "achievements": [
      "Reduced testing time by 60% through automation",
      "Mentored 5 junior QA engineers",
      "Implemented CI/CD testing pipeline"
    ]
  }'
```

**Expected Result:**
- Status: 201 Created
- Experience returned with generated ID
- `created_at` timestamp set

**Save Experience ID:**
```bash
export EXP_ID=$(curl -s -X POST http://localhost:8080/api/v1/profile/experience \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"company":"Tech Corp","position":"Senior QA Engineer","start_date":"2020-01-15","is_current":false}' \
  | jq -r '.id')
```

### Test Case 3.2: Create Experience - Current Position

**Objective:** Verify current position handling (no end_date)

```bash
curl -X POST http://localhost:8080/api/v1/profile/experience \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company": "Current Company Inc",
    "position": "Lead QA Engineer",
    "start_date": "2023-07-01",
    "is_current": true,
    "description": "Currently leading QA efforts"
  }'
```

**Expected Result:**
- Status: 201 Created
- `end_date` is null
- `is_current` is true

### Test Case 3.3: Create Experience - Missing Required Fields

**Objective:** Verify validation of required fields

```bash
curl -X POST http://localhost:8080/api/v1/profile/experience \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company": "Tech Corp"
  }'
```

**Expected Result:**
- Status: 400 Bad Request
- Error indicates missing required field (position)

### Test Case 3.4: Get All Experiences

**Objective:** Verify listing all user experiences

```bash
curl http://localhost:8080/api/v1/profile/experience \
  -H "Authorization: Bearer $QA_TOKEN"
```

**Expected Result:**
- Status: 200 OK
- Array of experiences (should have 2 from previous tests)
- Ordered by start_date or created_at

### Test Case 3.5: Update Experience

**Objective:** Verify experience update

```bash
curl -X PUT http://localhost:8080/api/v1/profile/experience/$EXP_ID \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company": "Tech Corp (Updated)",
    "position": "Staff QA Engineer",
    "start_date": "2020-01-15",
    "end_date": "2023-06-30",
    "is_current": false,
    "description": "Updated description",
    "achievements": [
      "Updated achievement 1",
      "Updated achievement 2"
    ]
  }'
```

**Expected Result:**
- Status: 200 OK
- Experience updated with new values

### Test Case 3.6: Update Non-existent Experience

**Objective:** Verify error handling for invalid ID

```bash
curl -X PUT http://localhost:8080/api/v1/profile/experience/00000000-0000-0000-0000-000000000000 \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company": "Test",
    "position": "Test"
  }'
```

**Expected Result:**
- Status: 404 Not Found
- Error code: `NOT_FOUND`

### Test Case 3.7: Delete Experience

**Objective:** Verify experience deletion

```bash
curl -X DELETE http://localhost:8080/api/v1/profile/experience/$EXP_ID \
  -H "Authorization: Bearer $QA_TOKEN"
```

**Expected Result:**
- Status: 200 OK
- Success message
- Experience no longer in database

**Verification:**
```bash
# Should return one less experience
curl http://localhost:8080/api/v1/profile/experience \
  -H "Authorization: Bearer $QA_TOKEN" | jq 'length'
```

---

## Education Management Testing

### Test Case 4.1: Create Education - Valid Data

**Objective:** Verify successful education creation

```bash
curl -X POST http://localhost:8080/api/v1/profile/education \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "institution": "State University",
    "degree": "Bachelor of Science",
    "field_of_study": "Computer Science",
    "start_date": "2014-09-01",
    "end_date": "2018-06-15",
    "gpa": 3.75
  }'
```

**Expected Result:**
- Status: 201 Created
- Education returned with generated ID

**Save Education ID:**
```bash
export EDU_ID=$(curl -s -X POST http://localhost:8080/api/v1/profile/education \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"institution":"State University","degree":"Bachelor of Science","field_of_study":"Computer Science","start_date":"2014-09-01","end_date":"2018-06-15"}' \
  | jq -r '.id')
```

### Test Case 4.2: Create Education - Without GPA

**Objective:** Verify GPA is optional

```bash
curl -X POST http://localhost:8080/api/v1/profile/education \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "institution": "Online University",
    "degree": "Master of Science",
    "field_of_study": "Software Engineering",
    "start_date": "2019-01-01",
    "end_date": "2021-12-31"
  }'
```

**Expected Result:**
- Status: 201 Created
- GPA is null

### Test Case 4.3: Create Education - Missing Required Fields

**Objective:** Verify validation of required fields

```bash
curl -X POST http://localhost:8080/api/v1/profile/education \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "institution": "University Name"
  }'
```

**Expected Result:**
- Status: 400 Bad Request
- Error indicates missing required fields (degree, field_of_study)

### Test Case 4.4: Get All Education

**Objective:** Verify listing all education entries

```bash
curl http://localhost:8080/api/v1/profile/education \
  -H "Authorization: Bearer $QA_TOKEN"
```

**Expected Result:**
- Status: 200 OK
- Array of education entries (should have 2 from previous tests)

### Test Case 4.5: Update Education

**Objective:** Verify education update

```bash
curl -X PUT http://localhost:8080/api/v1/profile/education/$EDU_ID \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "institution": "State University (Updated)",
    "degree": "Bachelor of Science",
    "field_of_study": "Computer Science",
    "start_date": "2014-09-01",
    "end_date": "2018-06-15",
    "gpa": 3.85
  }'
```

**Expected Result:**
- Status: 200 OK
- Education updated (GPA changed from 3.75 to 3.85)

### Test Case 4.6: Delete Education

**Objective:** Verify education deletion

```bash
curl -X DELETE http://localhost:8080/api/v1/profile/education/$EDU_ID \
  -H "Authorization: Bearer $QA_TOKEN"
```

**Expected Result:**
- Status: 200 OK
- Success message

---

## Skills Management Testing

### Test Case 5.1: Create Skill - Technical

**Objective:** Verify technical skill creation

```bash
curl -X POST http://localhost:8080/api/v1/profile/skills \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Go",
    "category": "technical",
    "proficiency_level": "expert"
  }'
```

**Expected Result:**
- Status: 201 Created
- Skill returned with ID

**Save Skill ID:**
```bash
export SKILL_ID=$(curl -s -X POST http://localhost:8080/api/v1/profile/skills \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Go","category":"technical","proficiency_level":"expert"}' \
  | jq -r '.id')
```

### Test Case 5.2: Create Skill - Different Categories

**Objective:** Verify all skill categories

```bash
# Technical skill
curl -X POST http://localhost:8080/api/v1/profile/skills \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"PostgreSQL","category":"technical","proficiency_level":"advanced"}'

# Soft skill
curl -X POST http://localhost:8080/api/v1/profile/skills \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Leadership","category":"soft","proficiency_level":"expert"}'

# Language skill
curl -X POST http://localhost:8080/api/v1/profile/skills \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"English","category":"language","proficiency_level":"expert"}'
```

**Expected Result:**
- All three skills created successfully
- Different categories stored correctly

### Test Case 5.3: Create Skill - Different Proficiency Levels

**Objective:** Verify all proficiency levels

```bash
curl -X POST http://localhost:8080/api/v1/profile/skills \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Python","category":"technical","proficiency_level":"beginner"}'

curl -X POST http://localhost:8080/api/v1/profile/skills \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Java","category":"technical","proficiency_level":"intermediate"}'
```

**Expected Result:**
- Both skills created
- Proficiency levels stored correctly

### Test Case 5.4: Create Skill - Missing Name

**Objective:** Verify name is required

```bash
curl -X POST http://localhost:8080/api/v1/profile/skills \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category": "technical",
    "proficiency_level": "expert"
  }'
```

**Expected Result:**
- Status: 400 Bad Request
- Error indicates name is required

### Test Case 5.5: Get All Skills

**Objective:** Verify listing all skills

```bash
curl http://localhost:8080/api/v1/profile/skills \
  -H "Authorization: Bearer $QA_TOKEN"
```

**Expected Result:**
- Status: 200 OK
- Array of all skills created (should have 6 total)
- Grouped by category if implemented

### Test Case 5.6: Delete Skill

**Objective:** Verify skill deletion

```bash
curl -X DELETE http://localhost:8080/api/v1/profile/skills/$SKILL_ID \
  -H "Authorization: Bearer $QA_TOKEN"
```

**Expected Result:**
- Status: 200 OK
- Success message

---

## Document Generation Testing

### Test Case 6.1: Generate Resume - Complete Profile

**Objective:** Verify resume generation with full profile data

**Prerequisites:** Ensure profile has experiences, education, and skills

```bash
curl -X POST http://localhost:8080/api/v1/generate/resume \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_description": "We are seeking a Senior QA Engineer with strong automation skills. You will lead testing efforts, mentor junior engineers, and establish best practices. Required: 5+ years QA experience, expertise in test automation, CI/CD knowledge.",
    "job_title": "Senior QA Engineer",
    "company_name": "Tech Innovators Inc",
    "template_id": "classic"
  }'
```

**Expected Result:**
- Status: 200 OK
- Document created with type: "resume"
- Content field populated with AI-generated resume
- `free_generations_left` decremented by 1
- Generation history record created with token counts and cost

**Verification:**
```bash
# Check remaining free generations
curl http://localhost:8080/api/v1/user/me \
  -H "Authorization: Bearer $QA_TOKEN" | jq '.free_generations_left'

# Check generation history
# (Use database query shown in Test Case 11.5)
```

### Test Case 6.2: Generate Resume - Minimal Job Description

**Objective:** Verify resume generation with minimal input

```bash
curl -X POST http://localhost:8080/api/v1/generate/resume \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_description": "Backend developer needed"
  }'
```

**Expected Result:**
- Status: 200 OK
- Resume generated even with minimal job description
- Title uses default format or job description

### Test Case 6.3: Generate Resume - Missing Job Description

**Objective:** Verify job_description is required

```bash
curl -X POST http://localhost:8080/api/v1/generate/resume \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_title": "Developer"
  }'
```

**Expected Result:**
- Status: 400 Bad Request
- Error indicates job_description is required

### Test Case 6.4: Generate Cover Letter - Complete Profile

**Objective:** Verify cover letter generation

```bash
curl -X POST http://localhost:8080/api/v1/generate/cover-letter \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_description": "We are looking for a passionate QA Engineer to join our growing team. You will work closely with developers to ensure quality across our products.",
    "job_title": "QA Engineer",
    "company_name": "Startup Ventures",
    "template_id": "modern"
  }'
```

**Expected Result:**
- Status: 200 OK
- Document created with type: "cover_letter"
- Content personalized to job description
- Free generation count decremented

### Test Case 6.5: Free Generation Limit - Exceed Limit

**Objective:** Verify free generation limit enforcement

```bash
# Check current free generations
FREE_GENS=$(curl -s http://localhost:8080/api/v1/user/me \
  -H "Authorization: Bearer $QA_TOKEN" | jq '.free_generations_left')

echo "Free generations left: $FREE_GENS"

# If FREE_GENS > 0, generate until depleted
# Then try one more generation

curl -X POST http://localhost:8080/api/v1/generate/resume \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"job_description":"Test job when quota exceeded"}'
```

**Expected Result (when quota = 0):**
- Status: 403 Forbidden
- Error code: `NO_FREE_GENERATIONS`
- Message: "No free generations left. Please upgrade to premium."

### Test Case 6.6: Premium User - Unlimited Generations

**Objective:** Verify premium users have unlimited generations

**Setup:**
```sql
-- Make user premium
UPDATE users SET is_premium = true
WHERE email = 'qa.test1@example.com';
```

```bash
# Generate multiple documents without limit
for i in {1..5}; do
  curl -s -X POST http://localhost:8080/api/v1/generate/resume \
    -H "Authorization: Bearer $QA_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"job_description\":\"Test job $i\"}" \
    | jq '.message'
done
```

**Expected Result:**
- All 5 generations succeed
- `free_generations_left` not decremented (stays at 0 or original value)

---

## Document Management Testing

### Test Case 7.1: Get All Documents

**Objective:** Verify listing all user documents

```bash
curl http://localhost:8080/api/v1/documents \
  -H "Authorization: Bearer $QA_TOKEN"
```

**Expected Result:**
- Status: 200 OK
- Array of all documents user has generated
- Each document has: id, type, title, template_id, status, created_at, updated_at

### Test Case 7.2: Get Specific Document

**Objective:** Verify retrieving single document

**Get a document ID first:**
```bash
DOC_ID=$(curl -s http://localhost:8080/api/v1/documents \
  -H "Authorization: Bearer $QA_TOKEN" | jq -r '.[0].id')

curl http://localhost:8080/api/v1/documents/$DOC_ID \
  -H "Authorization: Bearer $QA_TOKEN"
```

**Expected Result:**
- Status: 200 OK
- Full document details including content (JSONB)

### Test Case 7.3: Get Non-existent Document

**Objective:** Verify error handling for invalid document ID

```bash
curl http://localhost:8080/api/v1/documents/00000000-0000-0000-0000-000000000000 \
  -H "Authorization: Bearer $QA_TOKEN"
```

**Expected Result:**
- Status: 404 Not Found
- Error code: `NOT_FOUND`

### Test Case 7.4: Update Document

**Objective:** Verify document editing

```bash
curl -X PUT http://localhost:8080/api/v1/documents/$DOC_ID \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Resume Title",
    "content": {
      "raw_content": "This is manually edited content...",
      "custom_section": "Additional info"
    },
    "status": "final"
  }'
```

**Expected Result:**
- Status: 200 OK
- Document updated with new title, content, and status
- `updated_at` timestamp changed

### Test Case 7.5: Update Document - Change Status

**Objective:** Verify status transitions (draft → final)

```bash
curl -X PUT http://localhost:8080/api/v1/documents/$DOC_ID \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "final"
  }'
```

**Expected Result:**
- Status: 200 OK
- Only status changed to "final"

### Test Case 7.6: Delete Document

**Objective:** Verify document deletion

```bash
curl -X DELETE http://localhost:8080/api/v1/documents/$DOC_ID \
  -H "Authorization: Bearer $QA_TOKEN"
```

**Expected Result:**
- Status: 200 OK
- Success message
- Document removed from database

**Verification:**
```bash
# Should return 404
curl http://localhost:8080/api/v1/documents/$DOC_ID \
  -H "Authorization: Bearer $QA_TOKEN"
```

### Test Case 7.7: Delete Another User's Document

**Objective:** Verify users cannot delete other users' documents

**Setup:** Create second user and document, try to delete with first user's token

```bash
# This should fail with 404 or 403
curl -X DELETE http://localhost:8080/api/v1/documents/OTHER_USER_DOC_ID \
  -H "Authorization: Bearer $QA_TOKEN"
```

**Expected Result:**
- Status: 404 Not Found (document not accessible)
- OR Status: 403 Forbidden

---

## Authorization & Security Testing

### Test Case 8.1: Access Protected Route Without Token

**Objective:** Verify authentication requirement

```bash
curl http://localhost:8080/api/v1/profile
```

**Expected Result:**
- Status: 401 Unauthorized
- Error code: `UNAUTHORIZED` or similar
- Message about missing/invalid token

### Test Case 8.2: Access Protected Route With Invalid Token

**Objective:** Verify token validation

```bash
curl http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer invalid.token.here"
```

**Expected Result:**
- Status: 401 Unauthorized
- Error code: `INVALID_TOKEN`

### Test Case 8.3: Access Protected Route With Expired Token

**Objective:** Verify token expiration enforcement

**Setup:** Use a token that's older than JWT_ACCESS_EXPIRY (default 15m)

```bash
# Wait 15+ minutes or manually create expired token
curl http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer $EXPIRED_TOKEN"
```

**Expected Result:**
- Status: 401 Unauthorized
- Error about token expiration

### Test Case 8.4: User Isolation - Profile Access

**Objective:** Verify users can only access their own data

**Setup:**
```bash
# Create second user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"qa.test2@example.com","password":"SecurePass456!","full_name":"QA Test 2"}'

# Get second user's token
QA_TOKEN_2=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"qa.test2@example.com","password":"SecurePass456!"}' \
  | jq -r '.access_token')

# Add experience to user 1
EXP_USER1=$(curl -s -X POST http://localhost:8080/api/v1/profile/experience \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"company":"User1 Company","position":"Engineer","start_date":"2020-01-01","is_current":true}' \
  | jq -r '.id')

# Try to access user 1's experience with user 2's token
curl -X PUT http://localhost:8080/api/v1/profile/experience/$EXP_USER1 \
  -H "Authorization: Bearer $QA_TOKEN_2" \
  -H "Content-Type: application/json" \
  -d '{"company":"Hacked","position":"Hacker"}'
```

**Expected Result:**
- Status: 404 Not Found (experience not found for user 2)
- OR Status: 403 Forbidden
- User 2 cannot modify user 1's data

### Test Case 8.5: User Isolation - Documents

**Objective:** Verify document access isolation

```bash
# Create document for user 1
DOC_USER1=$(curl -s -X POST http://localhost:8080/api/v1/generate/resume \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"job_description":"Test isolation"}' \
  | jq -r '.document.id')

# Try to access with user 2's token
curl http://localhost:8080/api/v1/documents/$DOC_USER1 \
  -H "Authorization: Bearer $QA_TOKEN_2"
```

**Expected Result:**
- Status: 404 Not Found
- User 2 cannot see user 1's documents

### Test Case 8.6: SQL Injection Attempts

**Objective:** Verify protection against SQL injection

```bash
# Try SQL injection in email field
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com OR 1=1--","password":"anything"}'

# Try SQL injection in experience company field
curl -X POST http://localhost:8080/api/v1/profile/experience \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"company":"Test\"; DROP TABLE users; --","position":"Test","start_date":"2020-01-01"}'
```

**Expected Result:**
- No SQL errors
- Inputs treated as literal strings
- No database corruption

### Test Case 8.7: XSS Attempts in Profile Fields

**Objective:** Verify XSS protection

```bash
curl -X PUT http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "summary": "<script>alert(\"XSS\")</script>",
    "website_url": "javascript:alert(\"XSS\")"
  }'
```

**Expected Result:**
- Data stored safely
- Script tags not executed (frontend responsibility)
- Input sanitized or escaped

---

## Error Handling Testing

### Test Case 9.1: Malformed JSON Request

**Objective:** Verify handling of invalid JSON

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"pass"'  # Missing closing brace
```

**Expected Result:**
- Status: 400 Bad Request
- Error about malformed JSON

### Test Case 9.2: Missing Content-Type Header

**Objective:** Verify Content-Type validation

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -d '{"email":"test@test.com","password":"SecurePass123!","full_name":"Test"}'
```

**Expected Result:**
- Status: 400 Bad Request or 415 Unsupported Media Type
- OR request parsed successfully (depends on implementation)

### Test Case 9.3: Large Payload

**Objective:** Verify handling of oversized requests

```bash
# Generate large string for achievements
LARGE_ACHIEVEMENTS=$(python3 -c "import json; print(json.dumps(['Achievement ' + str(i) for i in range(10000)]))")

curl -X POST http://localhost:8080/api/v1/profile/experience \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"company\":\"Test\",\"position\":\"Test\",\"start_date\":\"2020-01-01\",\"achievements\":$LARGE_ACHIEVEMENTS}"
```

**Expected Result:**
- Status: 413 Payload Too Large OR 400 Bad Request
- OR request succeeds if within limits

### Test Case 9.4: Database Connection Loss

**Objective:** Verify error handling when database is unavailable

**Setup:**
```bash
# Stop database
make docker-down
```

```bash
curl http://localhost:8080/api/v1/health
```

**Expected Result:**
- Status: 503 Service Unavailable
- Error about database connectivity

**Cleanup:**
```bash
make docker-up
```

### Test Case 9.5: OpenAI API Failure

**Objective:** Verify handling of AI service failures

**Setup:** Temporarily set invalid OpenAI API key in .env and restart server

```bash
curl -X POST http://localhost:8080/api/v1/generate/resume \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"job_description":"Test with invalid API key"}'
```

**Expected Result:**
- Status: 500 Internal Server Error OR 502 Bad Gateway
- Error message about generation failure
- User's free_generations_left NOT decremented
- No document created

---

## Database Integrity Testing

### Test Case 10.1: Cascade Delete - User Deletion

**Objective:** Verify cascade deletes work correctly

```sql
-- Create test user with full profile
-- Then delete user and verify all related data is deleted

-- Get user ID
SELECT id FROM users WHERE email = 'qa.test1@example.com';

-- Count related records
SELECT
    (SELECT COUNT(*) FROM profiles WHERE user_id = 'USER_ID') as profile_count,
    (SELECT COUNT(*) FROM documents WHERE user_id = 'USER_ID') as doc_count,
    (SELECT COUNT(*) FROM generation_history WHERE user_id = 'USER_ID') as gen_count;

-- Delete user
DELETE FROM users WHERE email = 'qa.test1@example.com';

-- Verify all related data deleted
SELECT
    (SELECT COUNT(*) FROM profiles WHERE user_id = 'USER_ID') as profile_count,
    (SELECT COUNT(*) FROM documents WHERE user_id = 'USER_ID') as doc_count,
    (SELECT COUNT(*) FROM generation_history WHERE user_id = 'USER_ID') as gen_count;
-- All should be 0
```

### Test Case 10.2: Cascade Delete - Profile Deletion

**Objective:** Verify profile cascade deletes experiences/education/skills

```sql
-- Get profile ID
SELECT id FROM profiles WHERE user_id = (SELECT id FROM users WHERE email = 'qa.test2@example.com');

-- Count related records
SELECT
    (SELECT COUNT(*) FROM experiences WHERE profile_id = 'PROFILE_ID') as exp_count,
    (SELECT COUNT(*) FROM education WHERE profile_id = 'PROFILE_ID') as edu_count,
    (SELECT COUNT(*) FROM skills WHERE profile_id = 'PROFILE_ID') as skill_count;

-- Delete profile
DELETE FROM profiles WHERE id = 'PROFILE_ID';

-- Verify all related data deleted
SELECT
    (SELECT COUNT(*) FROM experiences WHERE profile_id = 'PROFILE_ID') as exp_count,
    (SELECT COUNT(*) FROM education WHERE profile_id = 'PROFILE_ID') as edu_count,
    (SELECT COUNT(*) FROM skills WHERE profile_id = 'PROFILE_ID') as skill_count;
-- All should be 0
```

### Test Case 10.3: Unique Constraint - Email

**Objective:** Verify email uniqueness is enforced at database level

```sql
-- Try to insert duplicate email directly
INSERT INTO users (id, email, password_hash, full_name, free_generations_left, is_premium)
VALUES (
    gen_random_uuid(),
    'qa.test1@example.com',  -- Duplicate email
    'hash',
    'Duplicate User',
    2,
    false
);
-- Should fail with unique constraint violation
```

### Test Case 10.4: Unique Constraint - User-Profile Relationship

**Objective:** Verify one profile per user

```sql
-- Try to create second profile for same user
INSERT INTO profiles (id, user_id)
VALUES (
    gen_random_uuid(),
    (SELECT id FROM users WHERE email = 'qa.test1@example.com')
);
-- Should fail with unique constraint violation
```

### Test Case 10.5: Updated_at Trigger

**Objective:** Verify auto-update triggers work

```sql
-- Check current updated_at
SELECT updated_at FROM users WHERE email = 'qa.test1@example.com';

-- Wait a second
SELECT pg_sleep(1);

-- Update user
UPDATE users SET full_name = 'Updated Name' WHERE email = 'qa.test1@example.com';

-- Verify updated_at changed
SELECT updated_at FROM users WHERE email = 'qa.test1@example.com';
-- Should be newer timestamp
```

---

## Performance Testing

### Test Case 11.1: Concurrent Registrations

**Objective:** Test handling of simultaneous registrations

```bash
# Run 10 concurrent registrations
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"perf.test$i@example.com\",\"password\":\"SecurePass123!\",\"full_name\":\"Perf Test $i\"}" &
done
wait

# Verify all created
curl -s http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"perf.test1@example.com","password":"SecurePass123!"}' | jq '.access_token'
```

**Expected Result:**
- All 10 users created successfully
- No deadlocks or race conditions
- Each has unique profile

### Test Case 11.2: Bulk Experience Creation

**Objective:** Test handling of many experiences for one user

```bash
# Create 50 experiences
for i in {1..50}; do
  curl -s -X POST http://localhost:8080/api/v1/profile/experience \
    -H "Authorization: Bearer $QA_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"company\":\"Company $i\",\"position\":\"Position $i\",\"start_date\":\"2020-01-01\",\"is_current\":false}" > /dev/null
done

# Verify all created and retrieval performance
time curl -s http://localhost:8080/api/v1/profile/experience \
  -H "Authorization: Bearer $QA_TOKEN" | jq 'length'
```

**Expected Result:**
- All 50 experiences created
- Retrieval completes in < 1 second
- Response size reasonable

### Test Case 11.3: Document Generation Performance

**Objective:** Measure time for document generation

```bash
# Time single generation
time curl -X POST http://localhost:8080/api/v1/generate/resume \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"job_description":"Performance test job description with moderate length to simulate real usage. Looking for experienced backend engineer."}'
```

**Expected Result:**
- Completes in 5-15 seconds (depends on OpenAI API)
- Consistent timing across multiple runs

### Test Case 11.4: Large Profile Data Generation

**Objective:** Test generation with maximum profile data

```bash
# Ensure user has:
# - Complete profile
# - 10+ experiences with achievements
# - 5+ education entries
# - 20+ skills

# Then generate resume
curl -X POST http://localhost:8080/api/v1/generate/resume \
  -H "Authorization: Bearer $QA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_description":"Detailed job description...",
    "job_title":"Senior Engineer",
    "company_name":"Tech Corp"
  }'
```

**Expected Result:**
- Generation succeeds
- All profile data incorporated into prompt
- No data truncation

### Test Case 11.5: Query Performance - Generation History

**Objective:** Verify efficient queries with large dataset

```sql
-- After creating many documents and generations
EXPLAIN ANALYZE
SELECT
    gh.created_at,
    d.type,
    d.title,
    gh.prompt_tokens,
    gh.completion_tokens,
    gh.total_cost
FROM generation_history gh
JOIN documents d ON d.id = gh.document_id
WHERE gh.user_id = 'USER_ID'
ORDER BY gh.created_at DESC
LIMIT 20;
-- Check execution time and index usage
```

**Expected Result:**
- Query uses indexes (user_id, created_at)
- Execution time < 50ms
- No sequential scans

---

## Bug Tracking Template

When a bug is found, document it using this template:

### Bug Report: [Short Description]

**Bug ID:** BUG-001

**Severity:** Critical / High / Medium / Low

**Component:** Authentication / Profile / Document Generation / etc.

**Test Case:** [Reference test case number]

**Description:**
[Detailed description of the bug]

**Steps to Reproduce:**
1. Step 1
2. Step 2
3. Step 3

**Expected Behavior:**
[What should happen]

**Actual Behavior:**
[What actually happens]

**Request:**
```bash
[Copy exact curl command]
```

**Response:**
```json
{
  "error": "..."
}
```

**Environment:**
- Server version: [commit hash or version]
- Database: PostgreSQL 15
- Go version: 1.24.0
- OpenAI model: gpt-4o-mini

**Logs:**
```
[Relevant server logs]
```

**Database State:**
```sql
-- Relevant queries showing database state
```

**Workaround:**
[If any workaround exists]

**Suggested Fix:**
[If you have suggestions]

**Status:** Open / In Progress / Resolved

---

## Summary Checklist

Use this checklist to track overall QA progress:

### Authentication
- [ ] User registration (valid, duplicate, weak password, invalid email)
- [ ] User login (valid, wrong password, non-existent user)
- [ ] Token refresh (valid, invalid token)
- [ ] Logout

### Profile Management
- [ ] Get empty profile
- [ ] Update profile (full, partial)
- [ ] Profile ownership verification

### Experience Management
- [ ] Create experience (valid, current position, missing fields)
- [ ] List experiences
- [ ] Update experience (valid, non-existent)
- [ ] Delete experience

### Education Management
- [ ] Create education (with GPA, without GPA, missing fields)
- [ ] List education
- [ ] Update education
- [ ] Delete education

### Skills Management
- [ ] Create skills (all categories, all proficiency levels)
- [ ] List skills
- [ ] Delete skills

### Document Generation
- [ ] Generate resume (complete profile, minimal data)
- [ ] Generate cover letter
- [ ] Free generation limit enforcement
- [ ] Premium user unlimited generations
- [ ] OpenAI API error handling

### Document Management
- [ ] List all documents
- [ ] Get specific document
- [ ] Update document (content, status)
- [ ] Delete document
- [ ] User isolation

### Security
- [ ] Auth required on protected routes
- [ ] Token validation
- [ ] Token expiration
- [ ] User data isolation
- [ ] SQL injection protection
- [ ] XSS protection

### Error Handling
- [ ] Malformed JSON
- [ ] Large payloads
- [ ] Database connection loss
- [ ] External API failures

### Database Integrity
- [ ] Cascade deletes (user, profile)
- [ ] Unique constraints (email, user-profile)
- [ ] Auto-update triggers

### Performance
- [ ] Concurrent operations
- [ ] Bulk data handling
- [ ] Query performance
- [ ] Document generation timing

---

## QA Automation Script

Complete end-to-end test script:

```bash
#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080/api/v1"
PASSED=0
FAILED=0

# Helper functions
pass() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
    ((PASSED++))
}

fail() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    ((FAILED++))
}

info() {
    echo -e "${YELLOW}ℹ INFO${NC}: $1"
}

# Test 1: Health Check
info "Test 1: Health Check"
HEALTH=$(curl -s $BASE_URL/health)
if echo "$HEALTH" | grep -q '"status":"ok"'; then
    pass "Health check successful"
else
    fail "Health check failed"
fi

# Test 2: User Registration
info "Test 2: User Registration"
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/auth/register \
    -H "Content-Type: application/json" \
    -d '{"email":"qa.auto@example.com","password":"SecurePass123!","full_name":"QA Automation"}')

if echo "$REGISTER_RESPONSE" | grep -q 'access_token'; then
    pass "User registration successful"
    TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.access_token')
else
    fail "User registration failed"
    exit 1
fi

# Test 3: Duplicate Registration
info "Test 3: Duplicate Email Registration"
DUP_RESPONSE=$(curl -s -w "%{http_code}" -X POST $BASE_URL/auth/register \
    -H "Content-Type: application/json" \
    -d '{"email":"qa.auto@example.com","password":"AnotherPass!","full_name":"Duplicate"}')

if echo "$DUP_RESPONSE" | grep -q '409'; then
    pass "Duplicate email rejected"
else
    fail "Duplicate email not rejected"
fi

# Test 4: Login
info "Test 4: User Login"
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"qa.auto@example.com","password":"SecurePass123!"}')

if echo "$LOGIN_RESPONSE" | grep -q 'access_token'; then
    pass "User login successful"
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token')
else
    fail "User login failed"
fi

# Test 5: Get Profile
info "Test 5: Get User Profile"
PROFILE=$(curl -s $BASE_URL/profile -H "Authorization: Bearer $TOKEN")
if echo "$PROFILE" | grep -q 'qa.auto@example.com'; then
    pass "Get profile successful"
else
    fail "Get profile failed"
fi

# Test 6: Update Profile
info "Test 6: Update Profile"
UPDATE_PROFILE=$(curl -s -X PUT $BASE_URL/profile \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"location":"Test City","summary":"QA automation test user"}')

if echo "$UPDATE_PROFILE" | grep -q 'Test City'; then
    pass "Update profile successful"
else
    fail "Update profile failed"
fi

# Test 7: Create Experience
info "Test 7: Create Experience"
EXP_RESPONSE=$(curl -s -X POST $BASE_URL/profile/experience \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"company":"Auto Test Corp","position":"QA Engineer","start_date":"2020-01-01","is_current":true}')

if echo "$EXP_RESPONSE" | grep -q 'Auto Test Corp'; then
    pass "Create experience successful"
    EXP_ID=$(echo "$EXP_RESPONSE" | jq -r '.id')
else
    fail "Create experience failed"
fi

# Test 8: Create Education
info "Test 8: Create Education"
EDU_RESPONSE=$(curl -s -X POST $BASE_URL/profile/education \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"institution":"Test University","degree":"Bachelor","field_of_study":"CS","start_date":"2014-09-01","end_date":"2018-06-01"}')

if echo "$EDU_RESPONSE" | grep -q 'Test University'; then
    pass "Create education successful"
else
    fail "Create education failed"
fi

# Test 9: Create Skills
info "Test 9: Create Skills"
SKILL_RESPONSE=$(curl -s -X POST $BASE_URL/profile/skills \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Automation","category":"technical","proficiency_level":"expert"}')

if echo "$SKILL_RESPONSE" | grep -q 'Automation'; then
    pass "Create skill successful"
else
    fail "Create skill failed"
fi

# Test 10: Generate Resume
info "Test 10: Generate Resume"
GEN_RESPONSE=$(curl -s -X POST $BASE_URL/generate/resume \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"job_description":"Looking for QA automation engineer","job_title":"QA Engineer","company_name":"Test Corp"}')

if echo "$GEN_RESPONSE" | grep -q 'generated successfully\|completed'; then
    pass "Generate resume successful"
    DOC_ID=$(echo "$GEN_RESPONSE" | jq -r '.document.id')
else
    fail "Generate resume failed"
    echo "$GEN_RESPONSE"
fi

# Test 11: Get Documents
info "Test 11: Get All Documents"
DOCS=$(curl -s $BASE_URL/documents -H "Authorization: Bearer $TOKEN")
DOC_COUNT=$(echo "$DOCS" | jq 'length')
if [ "$DOC_COUNT" -ge 1 ]; then
    pass "Get documents successful (count: $DOC_COUNT)"
else
    fail "Get documents failed"
fi

# Test 12: Access Without Token
info "Test 12: Unauthorized Access"
UNAUTH=$(curl -s -w "%{http_code}" $BASE_URL/profile)
if echo "$UNAUTH" | grep -q '401'; then
    pass "Unauthorized access properly rejected"
else
    fail "Unauthorized access not rejected"
fi

# Summary
echo ""
echo "=================================="
echo "QA Test Summary"
echo "=================================="
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo "Total: $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi
```

Save this as `qa-automation.sh`, make it executable, and run:

```bash
chmod +x qa-automation.sh
./qa-automation.sh
```

---

## Next Steps

1. **Manual Testing**: Execute all test cases in this document systematically
2. **Bug Documentation**: Use the bug tracking template for any issues found
3. **Automation**: Implement automated tests using Go's testing framework
4. **CI/CD Integration**: Add automated tests to CI pipeline
5. **Performance Baseline**: Establish performance benchmarks for critical operations
6. **Load Testing**: Use tools like Apache Bench or k6 for load testing
7. **Security Audit**: Conduct thorough security review
8. **Documentation**: Update API documentation based on QA findings

---

## Contact

For questions or issues with QA testing, contact the development team or file an issue in the project repository.
