export const build = (numberOfBricks) => {
    // get the body to later append each brick we create
    const body = document.querySelector('body');

    // count the bricks we create
    let brickCount = 0;

    const interval = setInterval(() => {
        if (brickCount >= numberOfBricks) {
            clearInterval(interval);
            return;
        }

        brickCount++;
        // create the brick in a div with the specified id
        const brick = document.createElement('div');
        brick.id = `brick-${brickCount}`;

        // determine if its a block in the middle
        const isMiddleCol = (brickCount - 1) % 3 === 1;

        if (isMiddleCol) {
            brick.dataset.foundation = 'true';
        }

        // append the brick at the body...
        // will take the class from css file .div
        body.append(brick);
        // timer goes outside the function we created but inside the setInterval
    }, 100);
};

export const repair = (...ids) => {
    ids.forEach(id => {
        const element = document.getElementById(id);
        if (element) {
            // check if the element is middle brick (foundation attribute)
            if (element.hasAttribute('data-foundation')) {
                element.dataset.repaired = 'in progress';
            } else {
                element.dataset.repaired = 'true';
            }
        }
    });
};

export const destroy = () => {
    // find all bricks
    const allBricks = document.querySelectorAll('div[id^="brick-"]');

    if (allBricks.length > 0) {
        const lastBrick = allBricks[allBricks.length - 1];
        lastBrick.remove();
    }
};