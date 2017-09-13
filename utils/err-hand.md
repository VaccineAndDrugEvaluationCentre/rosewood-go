type ErrorContext interface {
    ContainsError() bool
    SetError(err error)
    Error() error
}

Here is how you would use it inside the error-prone function:

func FragileFunction(ctx ErrorContext) {
    if ctx.ContainsError() {
        return
    }

    //some processing ...
    //if there is an error
    ctx.SetError(errors.New("some error"))  
    return
}


It allows you to do this:

ctx := NewErrorContext()

fridge := GetFridge(ctx)
egg := GetEgg(fridge, ctx) 
mixedEgg := MixAndSeason(egg, ctx)
friedEgg := Fry(mixedEgg, ctx)

if ctx.ContainsError() {
    fmt.Printf("failed to fry eggs: %s", ctx.Error().Error())
}

Or you can even do this:

ctxA := NewErrorContext()
ctxB := NewErrorContext()
ignored := NewErrorContext()

a := DoSomethingOne(ctxA)
b := DoSomethingTwo(a, ctxA)
c,d := DoSomethingThree(b, ctxB) //different context
if ctxB.ContainsError() {
   c = 1
   d = 2
}
e := DoSomethingFour(c, d, ctxA)
if ctxA.ContainsError() {
    fmt.Println("Failed To do A")
}

DoSomething(e, ignored) //error is ignored