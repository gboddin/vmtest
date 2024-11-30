.org $8000

STDOUT_DEVICE = $C000

reset:
    LDX #$00                   ;reset X
loop:
    LDA helloMessage,x         ;Load helloMessage+X in A
    INX                        ;Increment X
    CMP #$00                   ;Check for null byte in A
    BEQ exit                   ;If found branch to exit
    STA STDOUT_DEVICE          ;Store A at memory address used by output device
    JMP loop                   ;Go again

helloMessage:
   .asc "Hello world!"         ;Message
   .byte $0A                   ;\n
   .byte $00                   ;null

exit:
   .byte $00
