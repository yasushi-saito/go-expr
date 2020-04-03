%{
package expr

import (
  "text/scanner"
)

func lexpos(yylex yyLexer) scanner.Position {
  return yylex.(*parser).sc.Pos()
}

%}

%union {
  ast exprNode
  ident string
}

%start main

%token <ident> tokIdent
%token <ident> tokLetrec tokIn tokIf
%token <ast> tokLiteral
%token <ident> tokArrow tokEQ tokNEQ tokGE tokLE

%type<astlist> main toplevelExprList
%type<ast> expr toplevelExpr
%type<assign> binding
%type<assignlist> bindingList
%type<arglist> arglist

%nonassoc IFPREC
%left '-' '+'
%left '*' '/'

%%

main: expr { yylex.(*parser).result = $1 }

expr: tokLiteral
  | tokIdent '(' exprList ')' { $$ = &exprApply{Func:$1, Args:$3} }
  | expr '+' expr { $$ = &exprApply{pos: lexpos(yylex), Func: "builtin:+", Args: []exprNode{$1, $3})
  | expr '-' expr { $$ = &exprApply{pos: lexpos(yylex), Func: "builtin:+", Args: []exprNode{$1, $3})
  | '(' expr ')' { $$ = $2 }
