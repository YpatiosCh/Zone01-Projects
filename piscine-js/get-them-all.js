export const getArchitects = () => {
    // get all architects from 'a' tag
    const architects = Array.from(document.getElementsByTagName('a'));

    // get all non architects from 'span' tag (the rest)
    const nonArchitects = Array.from(document.getElementsByTagName('span'));

    return [architects, nonArchitects];
}

export const getClassical = () => {
    // get all architects first
    const [architects] = getArchitects(); 
    /* shorthand for:
        const result = getArchitects();
        const architects = result[0];
    */

    // now from those architects, get classical architects
    const classicalArchitects = architects.filter(architect => architect.classList.contains('classical'));

    // and then the non classical architects
    const nonClassicalArchitects = architects.filter(architect => !architect.classList.contains('classical'));

    return [classicalArchitects, nonClassicalArchitects];
}

export const getActive = () => {
    // get the classical architects
    const [classicalArchitects] = getClassical();

    // fitler the active classical architects 
    const activeArchitects = classicalArchitects.filter(architect => architect.classList.contains('active'));

    // filter inactive classical architects
    const inactiveArchitects = classicalArchitects.filter(architect => !architect.classList.contains('active'));

    return [activeArchitects, inactiveArchitects];
}

export const getBonannoPisano = () => {
    // get bonanno pisano from id
    const bonano = document.getElementById('BonannoPisano');

    // get all active architects and filter so that bonano is not in the list
    const [activeArchitects] = getActive();

    const remainingArchitects = activeArchitects.filter(architect => architect.id !== 'BonannoPisano')

    return [bonano, remainingArchitects];
}