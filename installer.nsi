Unicode true
!include "MUI2.nsh"
!include "FileFunc.nsh"
!include "LogicLib.nsh"

!define APPNAME       "Case Clicker"
!define COMPANY       "Foxyyyy"
!define VERSION       "1.3.1"
!define EXENAME       "Case Clicker.exe"
!define REGKEY        "Software\Microsoft\Windows\CurrentVersion\Uninstall\CaseClicker"

Name                  "${APPNAME}"
OutFile               "Case Clicker Setup.exe"
InstallDir            "$LOCALAPPDATA\Programs\${APPNAME}"
InstallDirRegKey      HKCU "Software\${APPNAME}" "InstallDir"
RequestExecutionLevel user
SetCompressor /SOLID lzma
BrandingText          "${APPNAME} ${VERSION}"

VIProductVersion      "1.3.1.0"
VIAddVersionKey       "ProductName"     "${APPNAME}"
VIAddVersionKey       "CompanyName"     "${COMPANY}"
VIAddVersionKey       "FileDescription" "${APPNAME} Setup"
VIAddVersionKey       "FileVersion"     "${VERSION}"
VIAddVersionKey       "ProductVersion"  "${VERSION}"
VIAddVersionKey       "LegalCopyright"  "${COMPANY}"

!define MUI_ICON   "app.ico"
!define MUI_UNICON "app.ico"
!define MUI_ABORTWARNING
!define MUI_WELCOMEFINISHPAGE_BITMAP "wizard.bmp"
!define MUI_UNWELCOMEFINISHPAGE_BITMAP "wizard.bmp"
!define MUI_HEADERIMAGE
!define MUI_HEADERIMAGE_BITMAP "header.bmp"

!define MUI_WELCOMEPAGE_TITLE "Install ${APPNAME}"
!define MUI_WELCOMEPAGE_TEXT  "An idle crate-opening clicker.$\r$\n$\r$\nClick crates, buy autoclicker bots, unbox skins and trade up for permanent bonuses.$\r$\n$\r$\nThis installs to your user folder — no admin rights needed.$\r$\n$\r$\nClick Next to continue."

!define MUI_FINISHPAGE_RUN "$INSTDIR\${EXENAME}"
!define MUI_FINISHPAGE_RUN_TEXT "Play ${APPNAME} now"
!define MUI_FINISHPAGE_TEXT "Your save file lives in %APPDATA%\CaseClicker and survives reinstalls."

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH

!insertmacro MUI_LANGUAGE "English"

Function .onInit
  ; A silent run means the running game launched us to update itself.
  ; Give its process a moment to exit so its files are no longer locked.
  IfSilent 0 +2
  Sleep 2500
FunctionEnd

Function .onInstSuccess
  ; After a silent (in-app) update, reopen the game so the swap feels seamless.
  IfSilent 0 +2
  Exec '"$INSTDIR\${EXENAME}"'
FunctionEnd

; ---------------------------------------------------------------- install
Section "${APPNAME} (required)" SecCore
  SectionIn RO
  SetOutPath "$INSTDIR"
  File "${EXENAME}"
  File "app.ico"

  ; update.txt points the app at its update feed. Never overwrite it on an
  ; upgrade - the user (or publisher) may have pointed it somewhere custom.
  IfFileExists "$INSTDIR\update.txt" +2 0
  File "update.txt"

  WriteRegStr HKCU "Software\${APPNAME}" "InstallDir" "$INSTDIR"
  WriteUninstaller "$INSTDIR\Uninstall.exe"

  ; Add/Remove Programs entry
  WriteRegStr   HKCU "${REGKEY}" "DisplayName"     "${APPNAME}"
  WriteRegStr   HKCU "${REGKEY}" "DisplayVersion"  "${VERSION}"
  WriteRegStr   HKCU "${REGKEY}" "Publisher"       "${COMPANY}"
  WriteRegStr   HKCU "${REGKEY}" "DisplayIcon"     "$INSTDIR\app.ico"
  WriteRegStr   HKCU "${REGKEY}" "UninstallString" '"$INSTDIR\Uninstall.exe"'
  WriteRegStr   HKCU "${REGKEY}" "QuietUninstallString" '"$INSTDIR\Uninstall.exe" /S'
  WriteRegStr   HKCU "${REGKEY}" "InstallLocation" "$INSTDIR"
  WriteRegDWORD HKCU "${REGKEY}" "NoModify" 1
  WriteRegDWORD HKCU "${REGKEY}" "NoRepair" 1
  ${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
  IntFmt $0 "0x%08X" $0
  WriteRegDWORD HKCU "${REGKEY}" "EstimatedSize" "$0"
SectionEnd

Section "Start Menu shortcut" SecStart
  CreateDirectory "$SMPROGRAMS\${APPNAME}"
  CreateShortcut "$SMPROGRAMS\${APPNAME}\${APPNAME}.lnk" "$INSTDIR\${EXENAME}" "" "$INSTDIR\app.ico" 0
  CreateShortcut "$SMPROGRAMS\${APPNAME}\Uninstall ${APPNAME}.lnk" "$INSTDIR\Uninstall.exe"
SectionEnd

Section "Desktop shortcut" SecDesk
  CreateShortcut "$DESKTOP\${APPNAME}.lnk" "$INSTDIR\${EXENAME}" "" "$INSTDIR\app.ico" 0
SectionEnd

LangString DESC_SecCore  ${LANG_ENGLISH} "The game itself. Around 6 MB."
LangString DESC_SecStart ${LANG_ENGLISH} "Add ${APPNAME} to your Start Menu."
LangString DESC_SecDesk  ${LANG_ENGLISH} "Put a shortcut on your Desktop."
!insertmacro MUI_FUNCTION_DESCRIPTION_BEGIN
  !insertmacro MUI_DESCRIPTION_TEXT ${SecCore}  $(DESC_SecCore)
  !insertmacro MUI_DESCRIPTION_TEXT ${SecStart} $(DESC_SecStart)
  !insertmacro MUI_DESCRIPTION_TEXT ${SecDesk}  $(DESC_SecDesk)
!insertmacro MUI_FUNCTION_DESCRIPTION_END

; ---------------------------------------------------------------- uninstall
Section "Uninstall"
  Delete "$INSTDIR\${EXENAME}"
  Delete "$INSTDIR\app.ico"
  Delete "$INSTDIR\update.txt"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir  "$INSTDIR"

  Delete "$SMPROGRAMS\${APPNAME}\${APPNAME}.lnk"
  Delete "$SMPROGRAMS\${APPNAME}\Uninstall ${APPNAME}.lnk"
  RMDir  "$SMPROGRAMS\${APPNAME}"
  Delete "$DESKTOP\${APPNAME}.lnk"

  DeleteRegKey HKCU "${REGKEY}"
  DeleteRegKey HKCU "Software\${APPNAME}"

  ; Offer to keep the save file rather than silently deleting progress
  IfSilent +6
  MessageBox MB_YESNO|MB_ICONQUESTION "Also delete your saved progress?$\r$\n$\r$\nChoose No to keep it for a future reinstall." IDNO +4
  Delete "$APPDATA\CaseClicker\save.json"
  Delete "$APPDATA\CaseClicker\save.json.tmp"
  RMDir  "$APPDATA\CaseClicker"
  Goto +1
SectionEnd
