lexer grammer SGoLexer;

program : functions;

functions : functions function
          | function;
function : FUNC ID '{' lines '}';

lines : lines line
      | line;

line : expr | stmt;

expr : make_expr
     | send_expr
     | recv_expr
     | close_expr
     | none_expr
     | go_expr;

make_expr : MAKE '(' ID ',' NUM ')';
close_expr : CLOSE ID;
send_expr : SEND ID;
recv_expr : RECV ID;
none_expr : NONE;
go_expr : GO ID;

stmt : for_stmt
     | if_stmt
     | select_stmt;

for_stmt : FOR '{' lines '}';
if_stmt : IF '{' lines '}' 
        (ELSE '{' lines '}')?;
select_stmt : SELECT '{' select_list '}';
select_list : select_list select_item
            | select_item;
select_item : (send_expr | recv_expr | DEFAULT) ':' '{' lines '}' ',';

FUNC : 'func' | 'FUNC';
MAKE : 'make' | 'MAKE';
CLOSE : 'close' | 'CLOSE';
SEND : 'send' | 'SEND';
RECV : 'recv' | 'RECV';
NONE : 'none' | 'NONE';
GO : 'go' | 'GO';
FOR : 'for' | 'FOR';
IF : 'if' | 'IF';
ELSE : 'else' | 'ELSE';
SELECT : 'select' | 'SELECT';
DEFAULT : 'default' | 'DEFAULT';

ID : [a-zA-Z_][a-zA-Z_0-9]*;
NUM : [0-9]+;




