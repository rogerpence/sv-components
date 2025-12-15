export interface Props {
    totalPages: number;
    pageNumber: number;
    navRoute: string;
}
declare const NextPrevAction: import("svelte").Component<Props, {}, "">;
type NextPrevAction = ReturnType<typeof NextPrevAction>;
export default NextPrevAction;
