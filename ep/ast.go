package ep

//
//// https://github.com/angular/angular/blob/master/packages/compiler/src/expression_parser/ast.ts
//
//type ParserError struct {
//	message string
//	input string
//	errLocation string
//}
//
//type ParseSpan struct {
//	start int
//	End int
//}
//
//func (p ParseSpan) toAbsolute(absoluteOffset int) AbsoluteSourceSpan {
//  return AbsoluteSourceSpan{absoluteOffset + p.start, absoluteOffset + p.End}
//}
//
//type AbsoluteSourceSpan struct {
//	start int
//	End int
//}
//
//type AST struct  {
//	span ParseSpan
//	sourceSpan AbsoluteSourceSpan
//}
//
//func (a AST) visit(visitor AstVisitor) {
//}
//
//func (a AST) toString() string {
//	return "AST"
//}
//
//type ASTWithName struct  {
//	AST
//	nameSpan AbsoluteSourceSpan
//}
//
//func (a ASTWithName) visit(visitor AstVisitor) {
//}
//
//func (a ASTWithName) toString() string {
//	return "AST"
//}
//
///**
// * Represents a quoted expression of the form:
// *
// * quote = prefix `:` uninterpretedExpression
// * prefix = identifier
// * uninterpretedExpression = arbitrary string
// *
// * A quoted expression is meant to be pre-processed by an AST transformer that
// * converts it into another AST that no longer contains quoted expressions.
// * It is meant to allow third-party developers to extend Angular template
// * expression language. The `uninterpretedExpression` part of the quote is
// * therefore not interpreted by the Angular's own expression parser.
// */
//
//type Quote struct  {
//	AST
//	prefix string
//	uninterpretedExpression string
//	location interface{}
//}
//
//func (a Quote) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.visitQuote(q, context)
//}
//
//func (a Quote) toString() string {
//	return "Quote"
//}
//
///////
//
//type EmptyExpr struct  {
//	AST
//}
//func (a EmptyExpr) visit(visitor AstVisitor, context interface{}) {
//
//}
//func (a EmptyExpr) toString() string {
//	return "AST"
//}
//
/////
//
//type ImplicitReceiver struct  {
//	AST
//}
//func (a ImplicitReceiver) visit(visitor AstVisitor, context interface{})interface{} {
//	return visitor.visitImplicitReceiver(a, context)
//}
//
//
/////
//
//type ThisReceiver struct  {
//	AST
//}
//func (a ThisReceiver) visit(visitor AstVisitor, context interface{})interface{} {
//	return visitor.visitThisReceiver(a, context)
//}
//
/////
//
//type Chain struct  {
//	AST
//	expressions []interface{}
//}
//func (a Chain) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.visitChain(a, context)
//}
//
/////
//
//type Conditional struct  {
//	AST
//	condition AST
//	trueExp AST
//	falseExp AST
//}
//func (a Conditional) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.visitConditional(a, context)
//}
//
/////
//
//type PropertyRead struct  {
//	ASTWithName
//	receiver AST
//	name string
//}
//func (a PropertyRead) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.visitPropertyRead(a, context)
//}
//
/////
//
//type PropertyWrite struct  {
//	ASTWithName
//	receiver AST
//	name string
//	value AST
//}
//func (a PropertyWrite) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.visitPropertyWrite(a, context)
//}
//
/////
//
//type SafePropertyRead struct  {
//	ASTWithName
//	receiver AST
//	name string
//}
//func (a SafePropertyRead) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.visitSafePropertyRead(a, context)
//}
//
/////
//
//type KeyedRead struct  {
//	AST
//	receiver AST
//	key string
//}
//func (a KeyedRead) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.visitKeyedRead(a, context)
//}
//
/////
//
//type SafeKeyedRead struct  {
//	AST
//	receiver AST
//	key string
//}
//func (a SafeKeyedRead) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.visitSafeKeyedRead(a, context)
//}
//
/////
//
//type KeyedWrite struct  {
//	AST
//	receiver AST
//	key string
//}
//func (a KeyedWrite) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.visitKeyedWrite(a, context)
//}
//
/////
//
//type BindingPipe struct  {
//	ASTWithName
//	exp AST
//	name string
//	args []interface{}
//}
//func (a BindingPipe) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.visitPipe(a, context)
//}
//
/////
//
//type LiteralPrimitive struct  {
//	AST
//	value interface{}
//}
//func (a LiteralPrimitive) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.visitLiteralPrimitive(a, context)
//}
//
/////
//
//type LiteralArray struct  {
//	AST
//	expressions []interface{}
//}
//func (a LiteralArray) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.visitLiteralArray(a, context)
//}
//
//type LiteralMapKey struct {
//	key string
//	quoted bool
//}
//
/////
//
//type LiteralMap struct  {
//	AST
//	keys []LiteralMapKey
//	values []interface{}
//}
//func (a LiteralMap) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.visitLiteralMap(a, context)
//}
//
/////
//
//type Interpolation struct  {
//	AST
//	public []interface{}
//	expressions []interface{}
//}
//func (a Interpolation) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.Interpolation(a, context)
//}
//
//
/////
//
//type Binary struct  {
//	AST
//	operation string
//	left AST
//	right AST
//}
//func (a Binary) visit(visitor AstVisitor, context interface{}) interface{} {
//	return visitor.visitBinary(a, context)
//}
//
/////
//
//type Unary struct  {
//	operation string
//	left AST
//	right AST
//	operator string
//	expr AST
//	binaryOp AST
//	binaryLeft AST
//	binaryRight AST
//}
//func (a Unary) visit(visitor AstVisitor, context interface{}) interface{} {
//	// TODO
//	//if (visitor.visitUnary != nil) {
//	//	return visitor.visitUnary(this, context);
//	//}
//
//	return visitor.visitBinary(this, context);
//
//	//return visitor.visitBinary(a, context)
//}
//
//func UnaryCreateMinus(span ParseSpan, sourceSpan AbsoluteSourceSpan, expr AST) Unary {
//	return Unary {
//
//	}
//}
