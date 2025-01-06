# hash-maker
 This is making hash windows app write by golang


## 사용방법 예시

### 일반 해시

```bash
.\hash-maker.exe -startPath C:\Users\User\study-flutter\update_test_app_1\build\windows\x64\runner\Release\.
```

flutter로 build windows했다면

`프로젝트명\build\windows\x64\runner\Release`

아래에 빌드 된다. 그렇기에

`프로젝트명\build\windows\x64\runner\Release\.` 이렇게 `\.` 이걸 추가해줘서 Release 안에를 hash해야한다.


### zip 해시

```bash
.\hash-maker.exe -zip -zipPath C:\Users\User\study-flutter\update_test_app_1\build\windows\x64\runner\Release.zip
```

이런식으로 zip파일 지정해서 해시를 뜰 수 있습니다.


## 사용 예시

```bash
# 기본 해시 생성
./hash-maker -startPath "C:\Project\build\Release"

# 특정 폴더를 ZIP으로 만들고 해시 생성
./hash-maker -zipfolder "C:\Project\build\Release" -zipname "update.zip"

# 기존 ZIP 파일에 해시 추가
./hash-maker -zip -zipPath "C:\Project\update.zip"
```